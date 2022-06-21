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
	"encoding/json"
	"log"
	"strings"

	"github.com/alibaba/ioc-golang/config"
	"github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/autowire"
)

type defaultTag struct {
}

var defaultTagParamLoaderSingleton autowire.ParamLoader

func GetDefaultTagParamLoader() autowire.ParamLoader {
	if defaultTagParamLoaderSingleton == nil {
		defaultTagParamLoaderSingleton = &defaultTag{}
	}
	return defaultTagParamLoaderSingleton
}

/*
Load support load param like:
```go
type Config struct {
	Address  string
	Password string
	DB       string
}
```

from field:

```go
NormalRedis  normalRedis.Redis  `normal:"Impl,address=127.0.0.1&password=xxx&db=0"`
```
*/
func (p *defaultTag) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || fi == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supported")
	}
	splitedTagValue := strings.Split(fi.TagValue, ",")
	if len(splitedTagValue) < 2 {
		return nil, errors.New("not supported")
	}
	kvs := strings.Split(splitedTagValue[1], "&")
	kvMaps := make(map[string]interface{})
	for _, kv := range kvs {
		splitedKV := strings.Split(kv, "=")
		if len(splitedKV) != 2 {
			return nil, errors.New("not supported")
		}

		expandValue, _ := config.ExpandConfigValueIfNecessary(splitedKV[1])
		kvMaps[splitedKV[0]] = expandValue
	}
	data, err := json.Marshal(kvMaps)
	if err != nil {
		log.Printf("error json marshal %s\n", err)
		return nil, err
	}
	param := sd.ParamFactory()
	if err := json.Unmarshal(data, param); err != nil {
		log.Printf("error jsonun marshal %s\n", err)
		return nil, err
	}
	return param, nil
}
