package inject

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func genIOCRPCClientStub(ctx *genall.GenerationContext, root *loader.Package, rpcServiceStructInfos []*markers.TypeInfo) {
	// apiRoot
	//fmt.Println("pkgpath =", root.PkgPath)
	loadedRoots, err := loader.LoadRoots(root.PkgPath + "/api")
	if err != nil {
		panic(err)
	}

	apiRoot := loadedRoots[0]

	//fmt.Println("loaded roots go files = ", apiRoot.GoFiles)
	apiRoot.NeedTypesInfo()

	imports := &importsList{
		byPath:  make(map[string]string),
		byAlias: make(map[string]string),
		pkg:     apiRoot,
	}
	// avoid confusing aliases by "reserving" the root package's name as an alias
	imports.byAlias[apiRoot.Name] = ""

	outContent := new(bytes.Buffer)

	// FIXME: gen rpc service client stub in ./api/
	for _, info := range rpcServiceStructInfos {
		c := &copyMethodMaker{
			pkg:         apiRoot,
			importsList: imports,
			codeWriter:  &codeWriter{out: outContent},
		}
		autowireAlias := c.NeedImport("github.com/alibaba/ioc-golang/autowire")
		rpcClientAlias := c.NeedImport("github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client")

		c.Line(`func init() {`)
		c.Linef(`%s.RegisterStructDescriptor(&%s.StructDescriptor{`, rpcClientAlias, autowireAlias)
		c.Linef(`Factory: func() interface{} {
			return &%sIOCRPCClient{}
		},`, info.Name)
		c.Line(`})`)
		c.Line(`}`)

		c.Linef("type %sIOCRPCClient struct {", info.Name)
		methods := parseMethodInfoFromGoFiles(info.Name, root.GoFiles)
		for _, m := range methods {
			importsAlias := m.GetImportAlias()
			if len(importsAlias) != 0 {
				for _, importAlias := range importsAlias {
					//fmt.Println("import alias ", importAlias)
					for _, rawFileImport := range info.RawFile.Imports {
						var originAlias string
						//fmt.Println(rawFileImport, splitedBySpace)
						if rawFileImport.Name != nil {
							originAlias = rawFileImport.Name.String()
						} else {
							splitedImport := strings.Split(rawFileImport.Path.Value, "/")
							originAlias = strings.TrimPrefix(splitedImport[len(splitedImport)-1], `"`)
							originAlias = strings.TrimSuffix(originAlias, `"`)
						}
						//fmt.Println("origin ", originAlias)
						if originAlias == importAlias {
							//fmt.Println(rawFileImport.Name) // todo get alias
							//fmt.Println("need import ", rawFileImport.Path.Value)
							toImport := strings.TrimPrefix(rawFileImport.Path.Value, `"`)
							toImport = strings.TrimSuffix(toImport, `"`)
							clientStubAlias := c.NeedImport(toImport)
							m.swapAlias(importAlias, clientStubAlias)
						}
					}
				}
			}
			c.Linef("%s func%s", m.name, m.body)
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
		outPath := filepath.Join(outAPIDir, fmt.Sprintf("/zz_generated.ioc_rpc_client_%s.go", strings.ToLower(info.Name)))
		file, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}

		writeOut(ctx, file, apiRoot, outBytes)
	}
}

type method struct {
	name string
	body string // like '(name, param *substruct.Param) (string, error)'
}

func getTailLetter(input string) string {
	i := len(input) - 1
	if i < 0 {
		return ""
	}
	for ; i >= 0; i-- {
		r := input[i]
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_') {
			return input[i+1:]
		}
	}
	return input
}

func (m *method) swapAlias(from, to string) {
	m.body = strings.Replace(m.body, from+".", to+".", -1)
}

func (m *method) GetImportAlias() []string {
	result := make([]string, 0)
	splitedByDot := strings.Split(m.body, ".")
	if len(splitedByDot) == 1 {
		return result
	}
	splitedByDotIgnoreFinal := splitedByDot[:len(splitedByDot)-1]

	for _, v := range splitedByDotIgnoreFinal {
		result = append(result, getTailLetter(v))
	}

	return result
}

/*
valid line is like
func (s *ServiceStruct) GetString(name string, param *substruct.Param) string {
func (s *ServiceStruct) GetString(param *substruct.Param) string {
func (s *ServiceStruct) GetString(name, param *substruct.Param) (string, error) {
func (s *ServiceStruct) GetString() string {
func (s *ServiceStruct) GetString()  {
*/
func newMethodFromLine(structName, line string) (method, bool) {
	line = strings.TrimSpace(line)
	funcPrefix := fmt.Sprintf("func (s *%s)", structName)
	if strings.HasPrefix(line, funcPrefix) && strings.HasSuffix(line, "{") {
		line = strings.TrimPrefix(line, funcPrefix)
		line = strings.TrimSuffix(line, "{")
		/*
			line can be
			GetString(param *substruct.Param) string
			GetString(name, param *substruct.Param) (string, error)
			GetString() string
			GetString()
		*/
		line = strings.TrimSpace(line)
		splited := strings.Split(line, "(")
		// funcName is 'GetString'
		funcName := strings.TrimSpace(splited[0])
		// funcBody is like '() string', '(name, param *substruct.Param) (string, error)'
		funcBody := strings.TrimSpace("(" + strings.Join(splited[1:], "("))
		return method{
			name: funcName,
			body: funcBody,
		}, true
	}
	return method{}, false
}

func parseMethodInfoFromGoFiles(structName string, goFilesPath []string) []method {
	allMethods := make([]method, 0)
	for _, filePath := range goFilesPath {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			//fmt.Printf("load file %s with error = %s\n", filePath, err)
			continue
		}
		fileString := string(data)
		fileLines := strings.Split(fileString, "\n")
		for _, line := range fileLines {
			parsedMethod, ok := newMethodFromLine(structName, line)
			if ok {
				allMethods = append(allMethods, parsedMethod)
			}
		}
	}
	return allMethods
}
