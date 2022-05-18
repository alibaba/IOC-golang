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
	"github.com/alibaba/ioc-golang/autowire"
)

type defaultParamLoader struct {
	defaultConfigParamLoader     autowire.ParamLoader
	defaultTagParamLoader        autowire.ParamLoader
	defaultTagPointToParamLoader autowire.ParamLoader
}

var defaultParamLoaderSingleton autowire.ParamLoader

func GetDefaultParamLoader() autowire.ParamLoader {
	if defaultParamLoaderSingleton == nil {
		defaultParamLoaderSingleton = &defaultParamLoader{
			defaultConfigParamLoader:     GetDefaultConfigParamLoader(),
			defaultTagParamLoader:        GetDefaultTagParamLoader(),
			defaultTagPointToParamLoader: GetDefaultTagPointToConfigParamLoader(),
		}
	}
	return defaultParamLoaderSingleton
}

/*
Load try to load config from 3 types: ordered from harsh to loose

1. Try to use defaultTagPointToParamLoader to load from field tag
2. Try to use defaultTagParamLoader to load from config
3. Try to use defaultConfigParamLoader to load from config pointed by tag

It will return with error if both way are failed.
```
*/
func (d *defaultParamLoader) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if param, err := d.defaultTagPointToParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warning
	if param, err := d.defaultTagParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warnin\

	return d.defaultConfigParamLoader.Load(sd, fi)
}
