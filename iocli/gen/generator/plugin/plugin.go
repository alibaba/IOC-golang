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

package plugin

import (
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type CodeGeneratorPluginForOneStruct interface {
	Name() string
	Type() Type

	Init(markers markers.MarkerValues)
	GenerateSDMetadataForOneStruct(c CodeWriter)
	GenerateInFileForOneStruct(c CodeWriter)
}

type Type int

const (
	AOP      = Type(1)
	Autowire = Type(2)
)

type CodeGeneratorPluginForPkg interface {
	GenerateCodeInPkg(c CodeWriter)
}
