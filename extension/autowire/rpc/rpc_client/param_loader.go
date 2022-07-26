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

package rpc_client

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/param_loader"
)

var defaultParamLoaderSingleton autowire.ParamLoader

func getDefaultParamLoader() autowire.ParamLoader {
	if defaultParamLoaderSingleton == nil {
		defaultParamLoaderSingleton = &paramLoader{
			defaultConfigParamLoader:     getConfigParamLoader(),
			defaultTagParamLoader:        param_loader.GetDefaultTagParamLoader(),
			defaultTagPointToParamLoader: getTagPointToConfigParamLoader(),
		}
	}
	return defaultParamLoaderSingleton
}

type paramLoader struct {
	defaultConfigParamLoader     autowire.ParamLoader
	defaultTagParamLoader        autowire.ParamLoader
	defaultTagPointToParamLoader autowire.ParamLoader
}

func (d *paramLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if param, err := d.defaultTagPointToParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warning
	if param, err := d.defaultTagParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warning

	return d.defaultConfigParamLoader.Load(sd, fi)
}
