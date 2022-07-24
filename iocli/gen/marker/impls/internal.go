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

package impls

import "sigs.k8s.io/controller-tools/pkg/markers"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type enableIOCGolangAutowireMarker struct {
}

func (m *enableIOCGolangAutowireMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire", markers.DescribesType, false))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireTypeMarker struct {
}

func (m *iocGolangAutowireTypeMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:type", markers.DescribesType, ""))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireParamMarker struct {
}

func (m *iocGolangAutowireParamMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:paramType", markers.DescribesType, ""))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireParamLoaderMarker struct {
}

func (m *iocGolangAutowireParamLoaderMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:paramLoader", markers.DescribesType, ""))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireConstructFuncMarker struct {
}

func (m *iocGolangAutowireConstructFuncMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:constructFunc", markers.DescribesType, ""))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireBaseTypeMarker struct {
}

func (m *iocGolangAutowireBaseTypeMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:baseType", markers.DescribesType, false))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireAliasMarker struct {
}

func (m *iocGolangAutowireAliasMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:alias", markers.DescribesType, ""))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireProxyMarker struct {
}

func (m *iocGolangAutowireProxyMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:proxy", markers.DescribesType, false))
}

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:proxy=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/iocli/gen/marker.DefinitionGetter

type iocGolangAutowireProxyAutoInjectionMarker struct {
}

func (m *iocGolangAutowireProxyAutoInjectionMarker) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:autowire:proxy:autoInjection", markers.DescribesType, false))
}
