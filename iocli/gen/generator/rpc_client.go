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

package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/extension/autowire/allimpls"
	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin/common"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func genIOCRPCClientStub(ctx *genall.GenerationContext, root *loader.Package, rpcServiceStructInfos []*markers.TypeInfo, debugMode bool) {
	// api folder root
	loadedRoots, err := loader.LoadRoots(root.PkgPath + "/api")
	if err != nil {
		panic(err)
	}

	apiRoot := loadedRoots[0]

	apiRoot.NeedTypesInfo()

	for _, info := range rpcServiceStructInfos {
		imports, err := GetimportsList(&importsListParam{
			pkg: apiRoot,
		})
		if err != nil {
			fmt.Printf("get import list error = %s\n", err)
			return
		}

		// 1. create all struct-level plugins
		allImplPluginsList, err := allimpls.GetImpl(util.GetSDIDByStructPtr(new(plugin.CodeGeneratorPluginForOneStruct)))
		if err != nil {
			panic(err)
		}
		// 2. sort plugins
		sort.Sort(plugin.CodeGeneratorPluginForOneStructSorter(allImplPluginsList.([]plugin.CodeGeneratorPluginForOneStruct)))
		allImplPlugins := allImplPluginsList.([]plugin.CodeGeneratorPluginForOneStruct)
		for _, p := range allImplPlugins {
			p.Init(*info)
		}

		// avoid confusing aliases by "reserving" the root package's name as an alias
		imports.byAlias[apiRoot.Name] = ""

		outContent := new(bytes.Buffer)
		c, err := GetcopyMethodMaker(&copyMethodMakerParam{
			pkg:         apiRoot,
			importsList: imports,
			outContent:  outContent,
		})
		if err != nil {
			fmt.Printf("get copy method maker error = %s\n", err)
			return
		}

		autowireAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire")
		normalAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire/normal")
		rpcClientAlias := c.NeedImport("github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client")

		c.Line(`func init() {`)
		// generate client stub factory
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, rpcClientAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%sIOCRPCClient{}
		},`, common.ToFirstCharLower(info.Name))

		// 5. gen metadata
		c.Line(`Metadata: map[string]interface{}{`)
		// 5.1 gen aop plugins metadata
		c.Line(`"aop": map[string]interface{}{`)
		for _, pluginImpl := range allImplPlugins {
			if pluginImpl.Type() == plugin.AOP {
				pluginImpl.GenerateSDMetadataForOneStruct(root, c)
			}
		}
		c.Line(`},`)

		// 5.2 gen autowire plugins metadata
		c.Line(`"autowire": map[string]interface{}{`)
		for _, pluginImpl := range allImplPlugins {
			if pluginImpl.Type() == plugin.Autowire {
				pluginImpl.GenerateSDMetadataForOneStruct(root, c)
			}
		}
		c.Line(`},`)
		c.Line(`},`)
		c.Line(`})`)

		// generate client stub proxy factory
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, normalAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%sIOCRPCClient_{}
		},`, common.ToFirstCharLower(info.Name))
		c.Line(`})`)
		c.Line(`}`)

		common.GenProxyStruct("IOCRPCClient_", c, []*markers.TypeInfo{info}, root, debugMode)
		common.GenInterface("IOCRPCClient", c, []*markers.TypeInfo{info}, root)

		c.Linef("type %sIOCRPCClient struct {", common.ToFirstCharLower(info.Name))
		methods := common.ParseExportedMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for idx := range methods {
			importsAlias := methods[idx].GetImportAlias()
			aliasSwapMap := make(map[string]string)
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
							aliasSwapMap[importAlias] = clientStubAlias
						}
					}
				}
				methods[idx].SwapAliasMap(aliasSwapMap)
			}
			c.Linef("%s func%s", methods[idx].Name, methods[idx].Body)
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
