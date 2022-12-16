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
	"strings"

	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
)

func GenInterface(interfaceSuffix string, c plugin.CodeWriter, needInterfaceStructInfos []*markers.TypeInfo, root *loader.Package) {
	for _, info := range needInterfaceStructInfos {
		// get all methods
		c.Linef(`type %s%s interface {`, info.Name, interfaceSuffix)
		methods := ParseExportedMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for idx := range methods {
			importsAlias := methods[idx].GetImportAlias()
			if len(importsAlias) != 0 {
				aliasSwapMap := make(map[string]string)
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
							aliasSwapMap[importAlias] = clientStubAlias
						}
					}
				}
				methods[idx].SwapAliasMap(aliasSwapMap)
			}
			c.Linef("%s %s", methods[idx].Name, methods[idx].Body)
		}
		c.Line("}")
		c.Line("")
	}
	c.Linef("")
}
