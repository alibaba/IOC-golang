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

package inject

import (
	"fmt"
	"go/types"
	"io"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"sigs.k8s.io/controller-tools/pkg/genall"

	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

const (
	PackagePathSeparator = "/"
	Dot                  = "."
)

type codeWriter struct {
	out io.Writer
}

// Line writes a single line.
func (c *codeWriter) Line(line string) {
	fmt.Fprintln(c.out, line)
}

// Linef writes a single line with formatting (as per fmt.Sprintf).
func (c *codeWriter) Linef(line string, args ...interface{}) {
	fmt.Fprintf(c.out, line+"\n", args...)
}

// importsList keeps track of required imports, automatically assigning aliases
// to import statement.
type importsList struct {
	byPath  map[string]string
	byAlias map[string]string

	pkg *loader.Package
}

// NeedImport marks that the given package is needed in the list of imports,
// returning the ident (import alias) that should be used to reference the package.
func (l *importsList) NeedImport(importPath string) string {
	// we get an actual path from Package, which might include venddored
	// packages if running on a package in vendor.
	if ind := strings.LastIndex(importPath, "/vendor/"); ind != -1 {
		importPath = importPath[ind+8: /* len("/vendor/") */]
	}

	// check to see if we've already assigned an alias, and just return that.
	alias, exists := l.byPath[importPath]
	if exists {
		return alias
	}

	// otherwise, calculate an import alias by joining path parts till we get something unique
	restPath, nextWord := path.Split(importPath)

	for otherPath, exists := "", true; exists && otherPath != importPath; otherPath, exists = l.byAlias[alias] {
		if restPath == "" {
			// do something else to disambiguate if we're run out of parts and
			// still have duplicates, somehow
			alias += "x"
		}

		// can't have a first digit, per Go identifier rules, so just skip them
		for firstRune, runeLen := utf8.DecodeRuneInString(nextWord); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(nextWord) {
			nextWord = nextWord[runeLen:]
		}

		// make a valid identifier by replacing "bad" characters with underscores
		nextWord = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return '_'
		}, nextWord)

		alias = nextWord + alias
		if len(restPath) > 0 {
			restPath, nextWord = path.Split(restPath[:len(restPath)-1] /* chop off final slash */)
		}
	}

	l.byPath[importPath] = alias
	l.byAlias[alias] = importPath
	return alias
}

// ImportSpecs returns a string form of each import spec
// (i.e. `alias "path/to/import").  Aliases are only present
// when they don't match the package name.
func (l *importsList) ImportSpecs() []string {
	res := make([]string, 0, len(l.byPath))
	for importPath, alias := range l.byPath {
		pkg := l.pkg.Imports()[importPath]
		if pkg != nil && pkg.Name == alias {
			// don't print if alias is the same as package name
			// (we've already taken care of duplicates).
			res = append(res, fmt.Sprintf("%q", importPath))
		} else {
			res = append(res, fmt.Sprintf("%s %q", alias, importPath))
		}
	}
	return res
}

// copyMethodMakers makes DeepCopy (and related) methods for Go types,
// writing them to its codeWriter.
type copyMethodMaker struct {
	pkg *loader.Package
	*importsList
	*codeWriter
}

type paramImplPair struct {
	paramName         string
	implName          string
	constructFuncName string
}

// GenerateMethodsFor makes init method
// for the given type, when appropriate
func (c *copyMethodMaker) GenerateMethodsFor(ctx *genall.GenerationContext, root *loader.Package, imports *importsList, infos []*markers.TypeInfo) {
	paramImplPairs := make([]paramImplPair, 0)
	rpcServiceStructInfos := make([]*markers.TypeInfo, 0)
	needProxyStructInfos := make([]*markers.TypeInfo, 0)
	getMethodGenerateCtxs := make([]getMethodGenerateCtx, 0)
	constructFunctionInfoNames := make([]string, 0)
	c.Line(`func init() {`)
	autowireAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire")
	for _, info := range infos {
		if len(info.Markers["ioc:autowire"]) == 0 {
			continue
		}
		if !info.Markers["ioc:autowire"][0].(bool) {
			continue
		}

		if len(info.Markers["ioc:autowire:type"]) == 0 {
			continue
		}
		autowireType := info.Markers["ioc:autowire:type"][0].(string)

		baseType := false
		if len(info.Markers["ioc:autowire:baseType"]) != 0 {
			baseType = info.Markers["ioc:autowire:baseType"][0].(bool)
		}
		paramType := ""

		if len(info.Markers["ioc:autowire:paramType"]) != 0 {
			paramType = info.Markers["ioc:autowire:paramType"][0].(string)
		}

		paramLoader := ""
		if len(info.Markers["ioc:autowire:paramLoader"]) != 0 {
			paramLoader = info.Markers["ioc:autowire:paramLoader"][0].(string)
		}

		constructFunc := ""
		if len(info.Markers["ioc:autowire:constructFunc"]) != 0 {
			constructFunc = info.Markers["ioc:autowire:constructFunc"][0].(string)
		}

		alise := ""
		if autowireType == "normal" || autowireType == "singleton" {
			alise = c.NeedImport(fmt.Sprintf("github.com/alibaba/ioc-golang/autowire/%s", autowireType))
		} else if autowireType == "rpc" {
			alise = c.NeedImport("github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_service")
			rpcServiceStructInfos = append(rpcServiceStructInfos, info)
		} else {
			alise = c.NeedImport(fmt.Sprintf("github.com/alibaba/ioc-golang/extension/autowire/%s", autowireType))
		}

		// gen proxy registry
		needProxyStructInfos = append(needProxyStructInfos, info)
		normalAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire/normal")
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, normalAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%s_{}
		},`, toFirstCharLower(info.Name))
		c.Line(`})`)

		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, alise, autowireAlias)

		// 0.gen alias
		if autowireType == "rpc" {
			c.Linef(`Alias: "%s/api.%sIOCRPCClient",`, root.PkgPath, info.Name)
		} else if len(info.Markers["ioc:autowire:alias"]) != 0 {
			c.Linef(`Alias: "%s",`, info.Markers["ioc:autowire:alias"][0].(string))
		}

		// 2. gen struct factory and gen param
		if !baseType {
			c.Linef(`Factory: func() interface{} {
			return &%s{}
		},`, info.Name)
			if paramType != "" {
				c.Line(`ParamFactory: func() interface{} {`)
				if constructFunc != "" && paramType != "" {
					c.Linef(`var _ %s = &%s{}`, getParamInterfaceType(paramType), paramType)
				}
				c.Linef(`return &%s{}
		},`, paramType)
			}
		} else {
			c.Linef(`Factory: func() interface{} {
			return new(%s)
		},`, info.Name)
			if paramType != "" {
				c.Linef(`ParamFactory: func() interface{} {
			return new(%s)
		},`, paramType)
			}
		}

		// 3. gen param loader
		if paramLoader != "" {
			c.Linef(`ParamLoader: &%s{},`, paramLoader)
		}

		// 4. gen constructor
		if constructFunc != "" && paramType != "" {
			paramImplPairs = append(paramImplPairs, paramImplPair{
				implName:          info.Name,
				paramName:         paramType,
				constructFuncName: constructFunc,
			})
			c.Linef(`ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(%s)
			impl := i.(*%s)
			return param.%s(impl)
		},`, getParamInterfaceType(paramType), info.Name, constructFunc)
		} else if constructFunc != "" && paramType == "" {
			// gen specific construct function without param

			c.Linef(`ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
	impl := i.(*%s)
	var constructFunc %sConstructFunc = %s
	return constructFunc(impl)
},`, info.Name, info.Name, constructFunc)
			constructFunctionInfoNames = append(constructFunctionInfoNames, info.Name)
		}
		c.Line(`})`)

		getMethodGenerateCtxs = append(getMethodGenerateCtxs, getMethodGenerateCtx{
			paramTypeName:     paramType,
			autowireTypeAlias: alise,
			structName:        info.Name,
			autowireType:      autowireType,
		})

		typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
		if typeInfo == types.Typ[types.Invalid] {
			root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type: %s", info.Name), info.RawSpec))
		}
	}
	c.Line(`}`)

	// gen param interface
	for _, paramImplPair := range paramImplPairs {
		c.Linef(`type %s interface {
			%s (impl *%s) (*%s,error)
		}`, getParamInterfaceType(paramImplPair.paramName), paramImplPair.constructFuncName, paramImplPair.implName, paramImplPair.implName)
	}

	// gen constructFunc signature
	for _, name := range constructFunctionInfoNames {
		c.Linef(`type %sConstructFunc func(impl *%s) (*%s, error)`, name, name, name)
	}

	// gen proxy struct
	genProxyStruct("_", c, needProxyStructInfos, root)

	// gen interface
	genInterface("IOCInterface", c, needProxyStructInfos, root)

	// gen get method and get interface method
	for _, g := range getMethodGenerateCtxs {
		if g.autowireType == "config" {
			continue
		}
		if g.paramTypeName != "" && g.autowireType != "rpc" {
			utilAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire/util")

			c.Linef(`func Get%s(p *%s)(*%s, error){
			i, err := %s.GetImpl(%s.GetSDIDByStructPtr(new(%s)), p)
			if err != nil {
				return nil, err
			}
			impl := i.(*%s)
			return impl, nil
		}`, g.structName, g.paramTypeName, g.structName, g.autowireTypeAlias, utilAlias, g.structName, g.structName)
			c.Line("")

			c.Linef(`func Get%sIOCInterface(p *%s)(%sIOCInterface, error){
				i, err := %s.GetImplWithProxy(%s.GetSDIDByStructPtr(new(%s)), p)
				if err != nil {
					return nil, err
				}
				impl := i.(%sIOCInterface)
				return impl, nil
			}`, g.structName, g.paramTypeName, g.structName, g.autowireTypeAlias, utilAlias, g.structName, g.structName)
			c.Line("")
		} else {
			utilAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire/util")

			c.Linef(`func Get%s()(*%s, error){`, g.structName, g.structName)
			if g.autowireType == "rpc" {
				c.Linef(`i, err := %s.GetImpl(%s.GetSDIDByStructPtr(new(%s)))`, g.autowireTypeAlias, utilAlias, g.structName)
			} else {
				c.Linef(`i, err := %s.GetImpl(%s.GetSDIDByStructPtr(new(%s)), nil)`, g.autowireTypeAlias, utilAlias, g.structName)
			}
			c.Linef(`if err != nil {
				return nil, err
			}
			impl := i.(*%s)
			return impl, nil
			}`, g.structName)
			c.Line("")

			c.Linef(`func Get%sIOCInterface()(%sIOCInterface, error){`, g.structName, g.structName)
			if g.autowireType == "rpc" {
				c.Linef(`i, err := %s.GetImplWithProxy(%s.GetSDIDByStructPtr(new(%s)))`, g.autowireTypeAlias, utilAlias, g.structName)
			} else {
				c.Linef(`i, err := %s.GetImplWithProxy(%s.GetSDIDByStructPtr(new(%s)), nil)`, g.autowireTypeAlias, utilAlias, g.structName)
			}
			c.Linef(`if err != nil {
				return nil, err
			}
			impl := i.(%sIOCInterface)
			return impl, nil
			}`, g.structName)
			c.Line("")
		}
	}

	// gen iocRPC client
	genIOCRPCClientStub(ctx, root, rpcServiceStructInfos)
}

func getParamInterfaceType(paramType string) string {
	return fmt.Sprintf("%sInterface", firstCharLower(paramType))
}

func firstCharLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

type getMethodGenerateCtx struct {
	paramTypeName     string
	autowireTypeAlias string
	structName        string
	autowireType      string
}

func toFirstCharLower(input string) string {
	return strings.ToLower(string(input[0])) + input[1:]
}

func genProxyStruct(proxySuffix string, c *copyMethodMaker, needProxyStructInfos []*markers.TypeInfo, root *loader.Package) {
	for _, info := range needProxyStructInfos {
		// get all methods
		c.Linef(`type %s%s struct {`, toFirstCharLower(info.Name), proxySuffix)
		methods := parseMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for idx := range methods {
			importsAlias := methods[idx].GetImportAlias()
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
							methods[idx].swapAlias(importAlias, clientStubAlias)
						}
					}
				}
			}
			c.Linef("%s_ func%s", methods[idx].name, methods[idx].body)
		}
		c.Line("}")
		c.Line("")

		for _, m := range methods {
			charDescriber := string(strings.ToLower(info.Name)[0])
			c.Linef(`func (%s *%s%s) %s%s{`, charDescriber, toFirstCharLower(info.Name), proxySuffix, m.name, m.body)
			if m.ReturnValueNum() > 0 {
				c.Linef(`return %s.%s_(%s)`, charDescriber, m.name, m.GetParamValues())
			} else {
				c.Linef(`%s.%s_(%s)`, charDescriber, m.name, m.GetParamValues())
			}
			c.Linef(`}`)
			c.Line("")
		}
	}
	c.Linef("")
}

func genInterface(interfaceSuffix string, c *copyMethodMaker, needInterfaceStructInfos []*markers.TypeInfo, root *loader.Package) {
	for _, info := range needInterfaceStructInfos {
		// get all methods
		c.Linef(`type %s%s interface {`, info.Name, interfaceSuffix)
		methods := parseMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for idx := range methods {
			importsAlias := methods[idx].GetImportAlias()
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
							methods[idx].swapAlias(importAlias, clientStubAlias)
						}
					}
				}
			}
			c.Linef("%s %s", methods[idx].name, methods[idx].body)
		}
		c.Line("}")
		c.Line("")
	}
	c.Linef("")
}
