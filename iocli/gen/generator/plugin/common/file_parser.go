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
	"io/ioutil"
	"strings"
)

/*
ParseMethodInfoFromGoFiles parse all methods, FIXME: now we don't support parse method signature with '\n' inner. like:
func (s *ComplexService) RPCBasicType(name string,
ageF64Ptr *float64) (string, error)
*/
func ParseMethodInfoFromGoFiles(structName string, goFilesPath []string) []method {
	allMethods := make([]method, 0)
	for _, filePath := range goFilesPath {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		fileString := string(data)
		fileLines := strings.Split(fileString, "\n")
		for _, line := range fileLines {
			parsedMethod, ok := newMethodFromLine(structName, line)
			if ok {
				allMethods = append(allMethods, parsedMethod)
			}
		}
	}
	return allMethods
}
