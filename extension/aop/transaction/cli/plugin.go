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

func (t *txCodeGenerationPlugin) Init(markers markers.MarkerValues) {
	for _, v := range markers[transactionFunctionAnnotation] {
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

func (t *txCodeGenerationPlugin) GenerateSDMetadataForOneStruct(c plugin.CodeWriter) {
	if len(t.txFunctionPairs) > 0 {
		c.Line(`"transaction": map[string]string{`)
		for _, pair := range t.txFunctionPairs {
			c.Linef(`"%s":"%s",`, pair.Name, pair.RollbackName)
		}
		c.Linef(`},`)
	}
}

func (t *txCodeGenerationPlugin) GenerateInFileForOneStruct(c plugin.CodeWriter) {
}

type txFunctionPair struct {
	Name         string
	RollbackName string
}
