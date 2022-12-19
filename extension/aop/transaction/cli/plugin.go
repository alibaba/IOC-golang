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

package cli

import (
	"strings"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin/common"

	"sigs.k8s.io/controller-tools/pkg/loader"

	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/extension/aop/transaction"
	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
)

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:implements=github.com/alibaba/ioc-golang/iocli/gen/generator/plugin.CodeGeneratorPluginForOneStruct
// +ioc:autowire:allimpls:autowireType=normal
// +ioc:autowire:constructFunc=create

type txCodeGenerationPlugin struct {
	txFunctionPairs []txFunctionPair
	structName      string
	info            markers.TypeInfo
}

func create(t *txCodeGenerationPlugin) (*txCodeGenerationPlugin, error) {
	t.txFunctionPairs = make([]txFunctionPair, 0)
	return t, nil
}

func (t *txCodeGenerationPlugin) Name() string {
	return transaction.Name
}

func (t *txCodeGenerationPlugin) Type() plugin.Type {
	return plugin.AOP
}

func (t *txCodeGenerationPlugin) Init(info markers.TypeInfo) {
	t.structName = info.Name
	t.info = info
	for _, v := range info.Markers[transactionFunctionAnnotation] {
		if txFuncMark, ok := v.(string); ok {
			txFuncPairRawStrings := strings.Split(txFuncMark, "-")
			if len(txFuncPairRawStrings) == 1 {
				t.txFunctionPairs = append(t.txFunctionPairs, txFunctionPair{
					Name: txFuncPairRawStrings[0],
				})
			} else if len(txFuncPairRawStrings) == 2 {
				t.txFunctionPairs = append(t.txFunctionPairs, txFunctionPair{
					Name:         txFuncPairRawStrings[0],
					RollbackName: txFuncPairRawStrings[1],
				})
			}
		}
	}
}

func (t *txCodeGenerationPlugin) GenerateSDMetadataForOneStruct(root *loader.Package, c plugin.CodeWriter) {
	if len(t.txFunctionPairs) > 0 {
		c.Line(`"transaction": map[string]string{`)
		for _, pair := range t.txFunctionPairs {
			c.Linef(`"%s":"%s",`, pair.Name, pair.RollbackName)
		}
		c.Linef(`},`)
	}
}

func (t *txCodeGenerationPlugin) GenerateInFileForOneStruct(root *loader.Package, c plugin.CodeWriter) {
	// check if needs parse
	needGenerate := false
	for _, pair := range t.txFunctionPairs {
		if pair.RollbackName == "" {
			continue
		}
		needGenerate = true
	}
	if !needGenerate {
		return
	}

	methods := common.ParseExportedMethodInfoFromGoFiles(t.structName, root.GoFiles)
	for idx := range methods {
		importsAlias := methods[idx].GetImportAlias()
		if len(importsAlias) != 0 {
			aliasSwapMap := make(map[string]string)
			for _, importAlias := range importsAlias {
				for _, rawFileImport := range t.info.RawFile.Imports {
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
						aliasSwapMap[importAlias] = clientStubAlias
					}
				}
			}
			methods[idx].SwapAliasMap(aliasSwapMap)
		}
	}

	methodStaticInfoMap := make(map[string]common.Method)
	for _, v := range methods {
		methodStaticInfoMap[v.Name] = v
	}

	for _, pair := range t.txFunctionPairs {
		if pair.RollbackName == "" {
			continue
		}
		c.Linef(`type %s%sTxFunction func (%s, errMsg string)`, common.ToFirstCharLower(t.structName), pair.Name, methodStaticInfoMap[pair.Name].ParamBody)
		c.Linef(`var _ %s%sTxFunction = (&%s{}).%s`, common.ToFirstCharLower(t.structName), pair.Name, t.structName, pair.RollbackName)
	}
}

type txFunctionPair struct {
	Name         string
	RollbackName string
}
