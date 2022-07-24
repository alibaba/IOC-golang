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
	"regexp"
	"strings"
)

type method struct {
	Name string
	Body string // like '(name, param *substruct.Param) (string, error)'

	params           []param
	returnValueTypes []string
	isVariadic       bool
}

func (m *method) parseParamAndReturnValues() {
	// split to params and values
	deep := 0
	tempStr := ""
	paramsAndvalue := make([]string, 0)
	for _, c := range m.Body {
		if c == '(' {
			deep += 1
			continue
		}
		if c == ')' {
			deep -= 1
			if deep == 0 {
				paramsAndvalue = append(paramsAndvalue, tempStr)
				tempStr = ""
			}
			continue
		}
		if deep <= 1 {
			tempStr += string(c)
		}
	}
	if tempStr != "" {
		paramsAndvalue = append(paramsAndvalue, tempStr)
	}

	m.params = parseParam(paramsAndvalue[0])
	if len(paramsAndvalue) > 1 {
		m.returnValueTypes = parseReturnValueTypes(paramsAndvalue[1])
	}

	if len(m.params) > 0 {
		finalParam := m.params[len(m.params)-1]
		m.isVariadic = strings.HasPrefix(finalParam.paramType, "...")
	}
}

func parseReturnValueTypes(input string) []string {
	if strings.TrimSpace(input) == "" {
		return make([]string, 0)
	}
	return strings.Split(input, ",")
}

func parseParam(input string) []param {
	params := make([]param, 0)
	if strings.TrimSpace(input) == "" {
		return params
	}

	parampairs := strings.Split(input, ",")

	for _, paramPair := range parampairs {
		allSplitedString := make([]string, 0)
		allNoneEmptySplitedString := make([]string, 0)

		items := strings.Split(paramPair, " ")
		for _, item := range items {
			allSplitedString = append(allSplitedString, strings.Split(item, "*")...)
		}
		for _, v := range allSplitedString {
			if v != "" {
				allNoneEmptySplitedString = append(allNoneEmptySplitedString, v)
			}
		}
		if len(allNoneEmptySplitedString) == 1 {
			params = append(params, param{
				paramType: "",
				value:     allNoneEmptySplitedString[0],
			})
		} else {
			params = append(params, param{
				paramType: allNoneEmptySplitedString[1],
				value:     allNoneEmptySplitedString[0],
			})
		}
	}
	return params
}

type param struct {
	value     string
	paramType string
}

// func (hs *Impl) RegisterRouterWithRawHttpHandler(path string, handler func(w http.ResponseWriter, r *http.Request), method string) {

func (m *method) GetParamValues() string {
	// remove func( xxx ) in body
	paramValues := make([]string, 0)
	for _, v := range m.params {
		paramValues = append(paramValues, v.value)
	}
	returnValues := strings.Join(paramValues, ",")
	if m.isVariadic {
		returnValues += "..."
	}
	return returnValues
}

func (m *method) ReturnValueNum() int {
	return len(m.returnValueTypes)
}

func getTailLetter(input string) string {
	i := len(input) - 1
	if i < 0 {
		return ""
	}
	for ; i >= 0; i-- {
		r := input[i]
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return input[i+1:]
		}
	}
	return input
}

func (m *method) SwapAliasMap(swapMap map[string]string) {
	splitedByDot := strings.Split(m.Body, ".")
	if len(splitedByDot) == 1 {
		return
	}
	splitedByDotIgnoreFinal := splitedByDot[:len(splitedByDot)-1]

	resultBody := ""

	for _, v := range splitedByDotIgnoreFinal {
		originAlias := getTailLetter(v)
		resultBody += strings.TrimSuffix(v, originAlias)
		swapedValue, ok := swapMap[originAlias]
		if ok {
			resultBody += swapedValue
		} else {
			resultBody += originAlias
		}
		resultBody += "."
	}

	resultBody += splitedByDot[len(splitedByDot)-1]
	m.Body = resultBody
	m.parseParamAndReturnValues()
}

func (m *method) GetImportAlias() []string {
	result := make([]string, 0)
	splitedByDot := strings.Split(m.Body, ".")
	if len(splitedByDot) == 1 {
		return result
	}

	aliasMap := make(map[string]bool)
	splitedByDotIgnoreFinal := splitedByDot[:len(splitedByDot)-1]

	for _, v := range splitedByDotIgnoreFinal {
		aliasMap[getTailLetter(v)] = true
	}
	for k := range aliasMap {
		result = append(result, k)
	}
	return result
}

/*
valid line is like
func (s *ServiceStruct) GetString(name string, param *substruct.Param) string {
func (s *ServiceStruct) GetString(param *substruct.Param) string {
func (s *ServiceStruct) GetString(name, param *substruct.Param) (string, error) {
func (s *ServiceStruct) GetString() string {
func (s *ServiceStruct) GetString()  {
*/
func newMethodFromLine(structName, line string) (method, bool) {
	line = strings.TrimSpace(line)
	if funcBody, ok := MatchFunctionByStructName(line, structName); ok {
		line = strings.TrimSuffix(funcBody, "{")
		/*
			line can be
			GetString(param *substruct.Param) string
			GetString(name, param *substruct.Param) (string, error)
			GetString() string
			GetString()
		*/
		line = strings.TrimSpace(line)
		splited := strings.Split(line, "(")
		// funcName is 'GetString'
		funcName := strings.TrimSpace(splited[0])
		// funcBody is like '() string', '(name, param *substruct.Param) (string, error)'
		funcBody := strings.TrimSpace("(" + strings.Join(splited[1:], "("))
		newMethod := method{
			Name: funcName,
			Body: funcBody,
		}
		newMethod.parseParamAndReturnValues()
		return newMethod, true
	}
	return method{}, false
}

func MatchFunctionByStructName(functionSignature, structName string) (string, bool) {
	splitedFunctionSignature := strings.Split(functionSignature, structName)
	if len(splitedFunctionSignature) <= 1 {
		return "", false
	}
	if !strings.HasPrefix(strings.TrimSpace(splitedFunctionSignature[1]), ")") {
		return "", false
	}
	// match func (
	regString := "^func\\(\\w+\\*$"
	signatureHeader := strings.Replace(splitedFunctionSignature[0], " ", "", -1)
	ok, err := regexp.MatchString(regString, signatureHeader)
	if err != nil || !ok {
		return "", false
	}
	return strings.TrimPrefix(strings.Join(splitedFunctionSignature[1:], structName), ")"), strings.HasSuffix(functionSignature, "{")
}
