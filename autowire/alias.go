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

package autowire

import (
	"fmt"
)

const (
	AliasKey = "alias"
)

// sdIDAliasMap a map of SDID alias-mapping named as alias container.
var sdIDAliasMap = make(map[string]string)

func registerAlias(alias, value string) {
	if _, ok := sdIDAliasMap[alias]; ok {
		panic(fmt.Sprintf("[Autowire Alias] Duplicate alias:[%s]", alias))
	}

	sdIDAliasMap[alias] = value
}

func GetSDIDByAliasIfNecessary(key string) string {
	if mappingSDID, ok := GetSDIDByAlias(key); ok {
		return mappingSDID
	}

	return key
}

func GetSDIDByAlias(alias string) (string, bool) {
	sdid, ok := sdIDAliasMap[alias]
	return sdid, ok
}
