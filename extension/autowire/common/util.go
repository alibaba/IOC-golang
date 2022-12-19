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

package common

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseInterfacePkgAndInterfaceName(interfaceID string) (string, string) {
	splited := strings.Split(interfaceID, ".")
	if len(splited) == 1 {
		return "", splited[0]
	}
	return strings.Join(splited[:len(splited)-1], "."), splited[len(splited)-1]
}

//IsExportedMethod  Is this an exported - upper case - name
func IsExportedMethod(name string) bool {
	s, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(s)
}
