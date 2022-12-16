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
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/autowire/common"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
)

const commonImplementsAnnotation = "ioc:autowire:implements"
const commonActiveProfileAnnotation = "ioc:autowire:activeProfile"
const commonLoadAtOnceAnnotation = "ioc:autowire:loadAtOnce"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:implements=github.com/alibaba/ioc-golang/iocli/gen/generator/plugin.CodeGeneratorPluginForOneStruct
// +ioc:autowire:allimpls:autowireType=normal
// +ioc:autowire:constructFunc=create

type commonCodeGenerationPlugin struct {
	implements    []string
	activeProfile string
	loadAtOnce    bool
	typeName      string
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

func (t *commonCodeGenerationPlugin) Init(info markers.TypeInfo) {
	t.typeName = info.Name
	markers := info.Markers
	for _, v := range markers[commonImplementsAnnotation] {
		t.implements = append(t.implements, v.(string))
	}

	activeProfile := ""
	if activeProfileValues := markers[commonActiveProfileAnnotation]; len(activeProfileValues) > 0 {
		activeProfile = activeProfileValues[0].(string)
	}
	t.activeProfile = activeProfile

	loadAtOnce := false

	if loadAtOnceValues := markers[commonLoadAtOnceAnnotation]; len(loadAtOnceValues) > 0 {
		loadAtOnce = loadAtOnceValues[0].(bool)
	}
	t.loadAtOnce = loadAtOnce
}

func (t *commonCodeGenerationPlugin) GenerateSDMetadataForOneStruct(root *loader.Package, w plugin.CodeWriter) {
	if len(t.implements) == 0 && !t.loadAtOnce {
		return
	}
	w.Linef(`"%s": map[string]interface{}{`, autowire.CommonMetadataKey)

	// generate implements and active profiles metadata
	if len(t.implements) > 0 {
		w.Linef(`"%s":[]interface{}{`, autowire.CommonImplementsMetadataKey)
		for _, interfaceID := range t.implements {
			interfacePkg, interfaceName := common.ParseInterfacePkgAndInterfaceName(interfaceID)
			if interfacePkg == "" || root.PkgPath == interfacePkg {
				w.Linef(`new(%s),`, interfaceName)
			} else {
				interfacePkgAlias := w.NeedImport(interfacePkg)
				w.Linef(`new(%s.%s),`, interfacePkgAlias, interfaceName)
			}
		}
		w.Linef(`},`)
		if t.activeProfile != "" {
			w.Linef(`"%s":"%s",`, autowire.CommonActiveProfileMetadataKey, t.activeProfile)
		}
	}

	// generate load at once metadata
	if t.loadAtOnce {
		w.Linef(`"%s":%t,`, autowire.CommonLoadAtOnceMetadataKey, true)
	}
	w.Line(`},`)
}

func (t *commonCodeGenerationPlugin) GenerateInFileForOneStruct(root *loader.Package, w plugin.CodeWriter) {
	// 1. get implements interface pkg and name
	for _, interfaceID := range t.implements {
		interfacePkg, interfaceName := common.ParseInterfacePkgAndInterfaceName(interfaceID)
		if interfacePkg == "" || root.PkgPath == interfacePkg {
			w.Linef("var _ %s = &%s{}", interfaceName, t.typeName)
		} else {
			interfacePkgAlias := w.NeedImport(interfacePkg)
			w.Linef("var _ %s.%s = &%s{}", interfacePkgAlias, interfaceName, t.typeName)
		}
	}
}
