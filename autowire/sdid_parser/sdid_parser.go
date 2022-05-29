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

package sdid_parser

import (
	"fmt"
	"strings"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
)

const (
	configTagKey           = "config"
	configExtensionPkgPath = "github.com/alibaba/ioc-golang/extension/config"
)

type defaultSDIDParser struct {
}

var defaultSDIDParserSingleton autowire.SDIDParser

func GetDefaultSDIDParser() autowire.SDIDParser {
	if defaultSDIDParserSingleton == nil {
		defaultSDIDParserSingleton = &defaultSDIDParser{}
	}
	return defaultSDIDParserSingleton
}

func (p *defaultSDIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	splitedTagValue := strings.Split(fi.TagValue, ",")
	interfaceName := fi.FieldType
	if interfaceName == "" {
		interfaceName = splitedTagValue[0]
	}
	// +ioc:autowire:alias=order
	injectStructName := splitedTagValue[0]
	if autowire.HasAlias(injectStructName) {
		return injectStructName, nil // by alias
	}

	// `config:"placeholder"`
	if fi.TagKey == configTagKey {
		// `config:"github.com/alibaba/ioc-golang/extension/config.ConfigInt,xxx.yyy.zzz"`
		if strings.HasPrefix(injectStructName, configExtensionPkgPath) {
			return injectStructName, nil
		}

		// `config:"ConfigString,xxx.yyy.zzz"`
		return fmt.Sprintf("%s.%s", configExtensionPkgPath, injectStructName), nil
	}

	interfaceFullName := interfaceName // InterfaceName
	// github.com/author/project/package/subPackage/targetPackage
	if len(fi.FieldTypePkgPath) > 0 {
		// github.com/author/project/package/subPackage/targetPackage.InterfaceName
		interfaceFullName = fi.FieldTypePkgPath + "." + interfaceName
	}
	injectStructFullName := injectStructName

	if interfaceFullName == injectStructFullName {
		return injectStructFullName, nil
	}

	return util.GetIdByNamePair(interfaceFullName, injectStructFullName), nil
}
