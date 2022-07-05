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
	"fmt"

	"github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/config"
)

type configParamLoader struct {
}

func getConfigParamLoaderPrefix(autowireType string, sd *autowire.StructDescriptor) string {
	structConfigPathKey := sd.Alias
	if structConfigPathKey == "" {
		structConfigPathKey, _ = util.ToRPCClientStubInterfaceSDID(util.GetSDIDByStructPtr(sd.Factory()))
	}
	return fmt.Sprintf("autowire%[1]s%[2]s%[1]s<%[3]s>%[1]sparam", config.YamlConfigSeparator, autowireType, structConfigPathKey)
}

var configParamLoaderSingleton autowire.ParamLoader

func getConfigParamLoader() autowire.ParamLoader {
	if configParamLoaderSingleton == nil {
		configParamLoaderSingleton = &configParamLoader{}
	}
	return configParamLoaderSingleton
}

func (p *configParamLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supporterd")
	}
	param := sd.ParamFactory()
	prefix := getConfigParamLoaderPrefix(fi.TagKey, sd)
	if err := config.LoadConfigByPrefix(prefix, param); err != nil {
		return nil, err
	}
	return param, nil
}
