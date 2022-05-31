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
	// fixme: out of range
	injectStructName := splitedTagValue[0]
	return autowire.GetSDIDByAliasIfNecessary(injectStructName), nil
}
