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
	"strings"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
)

const (
	configTagKey         = "config"
	packagePathSeparator = "/"
	dot                  = "."
	emptyString          = ""
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
	// +ioc:autowire:alias=xxx
	injectStructName := splitedTagValue[0]
	if autowire.HasAlias(injectStructName) {
		return injectStructName, nil // by alias
	}

	// handle plain struct bean, such as:
	// `config:"github.com/alibaba/ioc-golang/extension/config.ConfigInt,xxx.yyy.zzz"`
	if alias, ok := tryFindAlias(injectStructName); ok {
		return alias, nil
	}

	if interfaceName == injectStructName {
		return injectStructName, nil
	}

	return util.GetIdByNamePair(interfaceName, injectStructName), nil
}

func tryFindAlias(structName string) (string, bool) {
	if isEligibleInterfaceReferencePath(structName) {
		shortName := structName[strings.LastIndex(structName, ".")+1:]
		if autowire.HasAlias(shortName) {
			return shortName, true
		}
		if target, ok := tryOtherCase(shortName); ok {
			return target, true
		}
	}

	return structName, false
}

func isEligibleInterfaceReferencePath(interfaceReferencePath string) bool {
	return strings.Contains(interfaceReferencePath, packagePathSeparator) &&
		strings.LastIndex(interfaceReferencePath, dot) > 0 &&
		strings.LastIndex(interfaceReferencePath, dot) < len(interfaceReferencePath)-1 &&
		(strings.LastIndex(interfaceReferencePath, packagePathSeparator) < strings.LastIndex(interfaceReferencePath, dot))
}

func tryOtherCase(structName string) (string, bool) {
	return tryCamelCase(structName)
}

func tryCamelCase(structName string) (string, bool) {
	camelCaseName := util.ToCamelCase(structName)
	if autowire.HasAlias(camelCaseName) {
		return camelCaseName, true
	}
	snakeCase := util.ToSnakeCase(camelCaseName)
	if autowire.HasAlias(snakeCase) {
		return snakeCase, true
	}

	return structName, false
}
