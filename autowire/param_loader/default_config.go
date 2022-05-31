/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package param_loader

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/config"
)

type defaultConfig struct {
}

func getDefaultConfigPrefix(sd *autowire.StructDescriptor) string {
	structConfigPathKey := sd.Alias
	if structConfigPathKey == "" {
		structConfigPathKey = util.GetSDIDByStructPtr(sd.Factory())
	}
	return fmt.Sprintf("autowire%[1]s%[2]s%[1]s%[3]s%[1]sparam", config.YamlConfigSeparator, sd.AutowireType(), structConfigPathKey)
}

var defaultConfigParamLoaderSingleton autowire.ParamLoader

func GetDefaultConfigParamLoader() autowire.ParamLoader {
	if defaultConfigParamLoaderSingleton == nil {
		defaultConfigParamLoaderSingleton = &defaultConfig{}
	}
	return defaultConfigParamLoaderSingleton
}

/*
Load support load struct described like:
```go
normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory:   func() interface{}{
			return &Impl{}
		},
		ParamFactory: func() interface{}{
			return &Config{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			return i, nil
		},
	})
}

type Config struct {
	Address  string
	Password string
	DB       string
}
```
with
Autowire type 'normal'
StructName 'Impl'

from:

```yaml
autowire:
  normal:
      github.com/alibaba/ioc-golang/test.Impl:
        param:
          address: 127.0.0.1
          password: xxx
          db: 0
```
*/
func (p *defaultConfig) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supporterd")
	}
	param := sd.ParamFactory()
	prefix := getDefaultConfigPrefix(sd)
	if err := config.LoadConfigByPrefix(prefix, param); err != nil {
		return nil, err
	}
	return param, nil
}
