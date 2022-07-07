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

package inject

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func genIOCRPCClientStub(ctx *genall.GenerationContext, root *loader.Package, rpcServiceStructInfos []*markers.TypeInfo) {
	// api folder root
	loadedRoots, err := loader.LoadRoots(root.PkgPath + "/api")
	if err != nil {
		panic(err)
	}

	apiRoot := loadedRoots[0]

	apiRoot.NeedTypesInfo()

	for _, info := range rpcServiceStructInfos {
		imports := &importsList{
			byPath:  make(map[string]string),
			byAlias: make(map[string]string),
			pkg:     apiRoot,
		}
		// avoid confusing aliases by "reserving" the root package's name as an alias
		imports.byAlias[apiRoot.Name] = ""

		outContent := new(bytes.Buffer)
		c := &copyMethodMaker{
			pkg:         apiRoot,
			importsList: imports,
			codeWriter:  &codeWriter{out: outContent},
		}
		autowireAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire")
		normalAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire/normal")
		rpcClientAlias := c.NeedImport("github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client")

		// calculation tx pairs
		txFunctionPairs := make([]txFunctionPair, 0)
		for _, v := range info.Markers["ioc:tx:func"] {
			if txFuncMark, ok := v.(string); ok {
				txFuncPairRawStrings := strings.Split(txFuncMark, "-")
				if len(txFuncPairRawStrings) == 1 {
					txFunctionPairs = append(txFunctionPairs, txFunctionPair{
						Name: txFuncPairRawStrings[0],
					})
				} else if len(txFuncPairRawStrings) == 2 {
					txFunctionPairs = append(txFunctionPairs, txFunctionPair{
						Name:         txFuncPairRawStrings[0],
						RollbackName: txFuncPairRawStrings[1],
					})
				}
			}
		}

		c.Line(`func init() {`)
		// generate client stub factory
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, rpcClientAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%sIOCRPCClient{}
		},`, toFirstCharLower(info.Name))

		// generate TransactionMethodsMap
		if len(txFunctionPairs) > 0 {
			c.Linef(`TransactionMethodsMap: map[string]string{`)
			for _, pair := range txFunctionPairs {
				c.Linef(`"%s":"%s",`, pair.Name, pair.RollbackName)
			}
			c.Linef(`},`)
		}
		c.Line(`})`)

		// generate client stub proxy factory
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, normalAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%sIOCRPCClient_{}
		},`, toFirstCharLower(info.Name))
		c.Line(`})`)
		c.Line(`}`)

		genProxyStruct("IOCRPCClient_", c, []*markers.TypeInfo{info}, root)
		genInterface("IOCRPCClient", c, []*markers.TypeInfo{info}, root)

		c.Linef("type %sIOCRPCClient struct {", toFirstCharLower(info.Name))
		methods := parseMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for _, m := range methods {
			importsAlias := m.GetImportAlias()
			if len(importsAlias) != 0 {
				for _, importAlias := range importsAlias {
					for _, rawFileImport := range info.RawFile.Imports {
						var originAlias string
						if rawFileImport.Name != nil {
							originAlias = rawFileImport.Name.String()
						} else {
							splitedImport := strings.Split(rawFileImport.Path.Value, "/")
							originAlias = strings.TrimPrefix(splitedImport[len(splitedImport)-1], `"`)
							originAlias = strings.TrimSuffix(originAlias, `"`)
						}
						if originAlias == importAlias {
							toImport := strings.TrimPrefix(rawFileImport.Path.Value, `"`)
							toImport = strings.TrimSuffix(toImport, `"`)
							clientStubAlias := c.NeedImport(toImport)
							m.swapAlias(importAlias, clientStubAlias)
						}
					}
				}
			}
			c.Linef("%s func%s", m.name, m.body)
		}
		c.Line("}")

		outBytes := outContent.Bytes()

		outContent = new(bytes.Buffer)
		writeHeaderWithoutConstrain(root, outContent, "api", imports, "")
		writeMethods(root, outContent, outBytes)
		outBytes = outContent.Bytes()
		formattedBytes, err := format.Source(outBytes)
		if err != nil {
			apiRoot.AddError(err)
			// we still write the invalid source to disk to figure out what went wrong
		} else {
			outBytes = formattedBytes
		}

		// ensure the directory exists

		outAPIDir := filepath.Dir(root.CompiledGoFiles[0]) + "/api"
		if err := os.MkdirAll(outAPIDir, os.ModePerm); err != nil {
			panic(err)
		}
		outPath := filepath.Join(outAPIDir, fmt.Sprintf("zz_generated.ioc_rpc_client_%s.go", strings.ToLower(info.Name)))
		file, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}

		writeOut(ctx, file, apiRoot, outBytes)
	}
}

type method struct {
	name string
	body string // like '(name, param *substruct.Param) (string, error)'

	params           []param
	returnValueTypes []string
	isVariadic       bool
}

func (m *method) parseParamAndReturnValues() {
	// split to params and values
	deep := 0
	tempStr := ""
	paramsAndvalue := make([]string, 0)
	for _, c := range m.body {
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

func (m *method) swapAlias(from, to string) {
	m.body = strings.Replace(m.body, from+".", to+".", -1)
	m.parseParamAndReturnValues()
}

func (m *method) GetImportAlias() []string {
	result := make([]string, 0)
	splitedByDot := strings.Split(m.body, ".")
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
	if funcBody, ok := matchFunctionByStructName(line, structName); ok {
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
			name: funcName,
			body: funcBody,
		}
		newMethod.parseParamAndReturnValues()
		return newMethod, true
	}
	return method{}, false
}

func matchFunctionByStructName(functionSignature, structName string) (string, bool) {
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

/*
parseMethodInfoFromGoFiles parse all methods, FIXME: now we don't support parse method signature with '\n' inner. like:
func (s *ComplexService) RPCBasicType(name string, age int, age32 int32, age64 int64, ageF32 float32,
ageF64 float64, namePtr *string, agePtr *int, age32Ptr *int32, age64Ptr *int64, ageF32Ptr *float32,
ageF64Ptr *float64) (string, int, int32, int64, float32, float64, *string, *int, *int32, *int64, *float32, *float64)
*/
func parseMethodInfoFromGoFiles(structName string, goFilesPath []string) []method {
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
