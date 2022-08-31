package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin/common"
)

func genGoPluginAPIInterface(ctx *genall.GenerationContext, root *loader.Package, rpcServiceStructInfos []*markers.TypeInfo) {
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

		common.GenInterface("", c, []*markers.TypeInfo{info}, root)

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
		outPath := filepath.Join(outAPIDir, fmt.Sprintf("zz_generated.ioc_%s.go", strings.ToLower(info.Name)))
		file, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}

		writeOut(ctx, file, apiRoot, outBytes)
	}
}
