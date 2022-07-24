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

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:proxy=false
// +ioc:autowire:paramType=objectGenCtxParam
// +ioc:autowire:constructFunc=Init

// objectGenCtx contains the common info for generating deepcopy implementations.
// It mostly exists so that generating for a package can be easily tested without
// requiring a full set of output rules, etc.
type objectGenCtx struct {
	objectGenCtxParam
}

type objectGenCtxParam struct {
	Collector  *markers.Collector
	Checker    *loader.TypeChecker
	HeaderText string
}

func (o *objectGenCtxParam) Init(i *objectGenCtx) (*objectGenCtx, error) {
	i.objectGenCtxParam = *o
	return i, nil
}

// generateForPackage generates IOCGolang init and runtime.Object implementations for
// types in the given package, writing the formatted result to given writer.
// May return nil if source could not be generated.
func (ctx *objectGenCtx) generateForPackage(genCtx *genall.GenerationContext, root *loader.Package) []byte {
	ctx.Checker.Check(root)

	root.NeedTypesInfo()

	imports, err := GetimportsList(&importsListParam{
		pkg: root,
	})
	if err != nil {
		fmt.Printf("get import list error = %s\n", err)
		return nil
	}
	// avoid confusing aliases by "reserving" the root package's name as an alias
	imports.byAlias[root.Name] = ""

	infos := make([]*markers.TypeInfo, 0)
	if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		infos = append(infos, info)
	}); err != nil {
		root.AddError(err)
		return nil
	}
	outContent := new(bytes.Buffer)
	copyCtx, err := GetcopyMethodMaker(&copyMethodMakerParam{
		pkg:         root,
		importsList: imports,
		outContent:  outContent,
	})
	if err != nil {
		fmt.Printf("get copy method maker error = %s\n", err)
		return nil
	}

	needGen := false
	for _, info := range infos {
		if len(info.Markers["ioc:autowire"]) != 0 {
			needGen = true
			break
		}
	}
	if !needGen {
		return nil
	}

	copyCtx.generateMethodsFor(genCtx, root, imports, infos)

	outBytes := outContent.Bytes()

	outContent = new(bytes.Buffer)
	writeHeader(root, outContent, root.Name, imports, ctx.HeaderText)
	writeMethods(root, outContent, outBytes)

	outBytes = outContent.Bytes()
	formattedBytes, err := format.Source(outBytes)
	if err != nil {
		root.AddError(err)
		// we still write the invalid source to disk to figure out what went wrong
	} else {
		outBytes = formattedBytes
	}

	return outBytes
}
