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
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/autowire/common"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
)

const commonImplementsAnnotation = "ioc:autowire:implements"
const commonActiveProfileAnnotation = "ioc:autowire:activeProfile"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:implements=github.com/alibaba/ioc-golang/iocli/gen/generator/plugin.CodeGeneratorPluginForOneStruct
// +ioc:autowire:allimpls:autowireType=normal
// +ioc:autowire:constructFunc=create

type commonCodeGenerationPlugin struct {
	implements    []string
	activeProfile string
}

func create(t *commonCodeGenerationPlugin) (*commonCodeGenerationPlugin, error) {
	t.implements = make([]string, 0)
	return t, nil
}

func (t *commonCodeGenerationPlugin) Name() string {
	return autowire.MetadataKey
}

func (t *commonCodeGenerationPlugin) Type() plugin.Type {
	return plugin.Autowire
}

func (t *commonCodeGenerationPlugin) Init(markers markers.MarkerValues) {
	for _, v := range markers[commonImplementsAnnotation] {
		t.implements = append(t.implements, v.(string))
	}

	activeProfile := ""
	if activeProfileValues := markers[commonActiveProfileAnnotation]; len(activeProfileValues) > 0 {
		activeProfile = activeProfileValues[0].(string)
	}
	t.activeProfile = activeProfile
}

func (t *commonCodeGenerationPlugin) GenerateSDMetadataForOneStruct(w plugin.CodeWriter) {
	if len(t.implements) > 0 {
		w.Linef(`"%s": map[string]interface{}{
				"%s":[]interface{}{`, autowire.CommonMetadataKey, autowire.CommonImplementsMetadataKey)
		for _, interfaceID := range t.implements {
			interfacePkg, interfaceName := common.ParseInterfacePkgAndInterfaceName(interfaceID)
			interfacePkgAlias := w.NeedImport(interfacePkg)
			w.Linef(`new(%s.%s),`, interfacePkgAlias, interfaceName)
		}
		w.Linef(`},`)
		if t.activeProfile != "" {
			w.Linef(`"%s":"%s",`, autowire.CommonActiveProfileMetadataKey, t.activeProfile)
		}
		w.Line(`},`)
	}
}

func (t *commonCodeGenerationPlugin) GenerateInFileForOneStruct(w plugin.CodeWriter) {
}
