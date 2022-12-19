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

	"github.com/alibaba/ioc-golang/extension/autowire/common"
)

/*
ParseExportedMethodInfoFromGoFiles parse all Upper case methods,
*/
func ParseExportedMethodInfoFromGoFiles(structName string, goFilesPath []string) []Method {
	exportedMethods := make([]Method, 0)
	for _, filePath := range goFilesPath {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		fileString := string(data)
		fileLines := strings.Split(fileString, "\n")
		fileLines = joinMethodLine(fileLines)
		for _, line := range fileLines {
			parsedMethod, ok := newMethodFromLine(structName, line)
			if ok {
				if common.IsExportedMethod(parsedMethod.Name) {
					exportedMethods = append(exportedMethods, parsedMethod)
				}
			}
		}
	}
	return exportedMethods
}

/*
joinMethodLine join func line splited by '\n' like:

func (s *ComplexService) RPCBasicType(name string,
	ageF64Ptr *float64,
) (string, error){
	return "", nil
}

after join, the output is
func (s *ComplexService) RPCBasicType(name string, ageF64Ptr *float64) (string, error){
	return "", nil
}
*/
func joinMethodLine(splitedFileLines []string) []string {
	afterJoinedFileLine := make([]string, 0)
	var inMethodSignature = false
	methodSignatureLine := make([]string, 0)
	for _, l := range splitedFileLines {
		l = validateSignatureSingleLine(l)
		if !inMethodSignature && strings.HasPrefix(strings.TrimSpace(l), "func ") && !strings.HasSuffix(strings.TrimSpace(l), "{") {
			inMethodSignature = true
			methodSignatureLine = append(methodSignatureLine, strings.TrimSpace(l))
		} else if inMethodSignature && strings.HasSuffix(strings.TrimSpace(l), "{") {
			inMethodSignature = false
			methodSignatureLine = append(methodSignatureLine, strings.TrimSpace(l))
		} else if inMethodSignature {
			methodSignatureLine = append(methodSignatureLine, strings.TrimSpace(l))
		}

		if !inMethodSignature {
			if len(methodSignatureLine) > 0 {
				// join signature lines
				joinedSignatoreLine := strings.Join(methodSignatureLine, " ")
				afterJoinedFileLine = append(afterJoinedFileLine, validateSignatureLine(joinedSignatoreLine), "\n")
				methodSignatureLine = make([]string, 0)
			} else {
				afterJoinedFileLine = append(afterJoinedFileLine, l, "\n")
			}
		}
	}
	return afterJoinedFileLine
}

// validateSignatureLine can validate joined signature line, replace ",)" to ")"
func validateSignatureLine(signatureLine string) string {
	splited := strings.Split(signatureLine, ")")
	for idx := range splited {
		splited[idx] = strings.TrimSpace(splited[idx])
		splited[idx] = strings.TrimSuffix(splited[idx], ",")
	}
	return strings.Join(splited, ")")
}

// validateSignatureSingleLine remote comments
func validateSignatureSingleLine(signatureLine string) string {
	splited := strings.Split(signatureLine, "//")
	return strings.TrimSpace(splited[0])
}
