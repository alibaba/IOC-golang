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
	"go/ast"

	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/extension/autowire/allimpls"

	"github.com/alibaba/ioc-golang/iocli/gen/marker"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var DebugMode = false

type Generator struct {
	HeaderFile string `marker:",optional"`
	Year       string `marker:",optional"`
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
}

// RegisterMarkers is called in main, register all markers
func (Generator) RegisterMarkers(into *markers.Registry) error {
	allImplGettersVal, err := allimpls.GetImpl(util.GetSDIDByStructPtr(new(marker.DefinitionGetter)))
	if err != nil {
		return err
	}
	defs := make([]*markers.Definition, 0)
	for _, g := range allImplGettersVal.([]marker.DefinitionGetter) {
		defs = append(defs, g.GetMarkerDefinition())
	}
	return markers.RegisterAll(into, defs...)
}

func (d Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string

	if d.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(d.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}

	// 1. get generate context
	objGenCtx, err := GetobjectGenCtx(&objectGenCtxParam{
		Collector:  ctx.Collector,
		Checker:    ctx.Checker,
		HeaderText: headerText,
		DebugMode:  DebugMode,
	})
	if err != nil {
		return err
	}

	for _, root := range ctx.Roots {
		// 2. generate codes under current pkg
		outContents := objGenCtx.generateForPackage(ctx, root)
		if outContents == nil {
			continue
		}
		// 3. write codes to file
		writeOut(ctx, nil, root, outContents)
	}
	return nil
}
