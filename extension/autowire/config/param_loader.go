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

package config

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
)

type paramLoader struct {
}

/*
Load support load config field like:
```go
Address  configString.ConfigString `config:"ConfigString,myConfig.myConfigSubPath.myConfigKey"`
```go

from:

```yaml
myConfig:
  myConfigSubPath:
      myConfigKey: myConfigValue
```
*/
func (p *paramLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || fi == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supported")
	}
	splitedTagValue := strings.Split(fi.TagValue, ",")
	configTagValue := splitedTagValue[1]
	param := sd.ParamFactory()
	if strings.HasPrefix(configTagValue, config.EnvPrefixKey) && strings.HasSuffix(configTagValue, config.EnvSuffixKey) {
		// config is env var
		configEnvVal := os.Getenv(strings.TrimSuffix(strings.TrimPrefix(configTagValue, config.EnvPrefixKey), config.EnvSuffixKey))
		if configEnvVal == "" {
			logrus.Warningf("load env key %s with empty string", configTagValue)
		}
		_ = yaml.Unmarshal([]byte(configEnvVal), param)
		return param, nil
	}
	// config is config file path
	if err := config.LoadConfigByPrefix(configTagValue, param); err != nil {
		logrus.Errorf("load config path %s error = %s", configTagValue, err.Error())
		// FIXME ignore config read error?
	}
	return param, nil
}

type Config struct {
	Address string
}
