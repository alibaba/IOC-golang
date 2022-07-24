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

	"github.com/alibaba/ioc-golang/extension/autowire/allimpls"
	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
)

const allimplsInterfaceAnnotation = "ioc:autowire:allimpls:interface"
const allimplsAutowireTypeAnnotation = "ioc:autowire:allimpls:autowireType"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/generator/plugin.CodeGeneratorPluginForOneStruct
// +ioc:autowire:allimpls:autowireType=normal
// +ioc:autowire:constructFunc=create

type allImplsCodeGenerationPlugin struct {
	implInterfaceIDs     []string
	allimplsAutowireType string
}

func create(t *allImplsCodeGenerationPlugin) (*allImplsCodeGenerationPlugin, error) {
	t.implInterfaceIDs = make([]string, 0)
	return t, nil
}

func (t *allImplsCodeGenerationPlugin) Name() string {
	return allimpls.Name
}

func (t *allImplsCodeGenerationPlugin) Type() plugin.Type {
	return plugin.Autowire
}

func (t *allImplsCodeGenerationPlugin) Init(markers markers.MarkerValues) {
	for _, v := range markers[allimplsInterfaceAnnotation] {
		t.implInterfaceIDs = append(t.implInterfaceIDs, v.(string))
	}

	allimplsAutowireType := ""
	if allimplsAutowireTypeValues := markers[allimplsAutowireTypeAnnotation]; len(allimplsAutowireTypeValues) > 0 {
		allimplsAutowireType = allimplsAutowireTypeValues[0].(string)
	}
	t.allimplsAutowireType = allimplsAutowireType
}

func (t *allImplsCodeGenerationPlugin) GenerateSDMetadataForOneStruct(c plugin.CodeWriter) {
	if len(t.implInterfaceIDs) > 0 {
		c.Linef(`"%s": map[string]interface{}{
				"%s":[]interface{}{`, allimpls.Name, allimpls.InterfaceMetadataKey)
		for _, interfaceID := range t.implInterfaceIDs {
			interfacePkg, interfaceName := parseInterfacePkgAndInterfaceName(interfaceID)
			interfacePkgAlias := c.NeedImport(interfacePkg)
			c.Linef(`new(%s.%s),`, interfacePkgAlias, interfaceName)
		}
		c.Linef(`},`)
		if t.allimplsAutowireType != "" {
			c.Linef(`"%s":"%s",`, allimpls.AutowireTypeMetadataKey, t.allimplsAutowireType)
		}
		c.Line(`},`)
	}
}

func (t *allImplsCodeGenerationPlugin) GenerateInFileForOneStruct(c plugin.CodeWriter) {
}

func parseInterfacePkgAndInterfaceName(interfaceID string) (string, string) {
	splited := strings.Split(interfaceID, ".")
	return strings.Join(splited[:len(splited)-1], "."), splited[len(splited)-1]
}
