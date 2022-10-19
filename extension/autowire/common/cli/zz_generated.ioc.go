//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package cli

import (
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	allimpls "github.com/alibaba/ioc-golang/extension/autowire/allimpls"
	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin"
	marker "github.com/alibaba/ioc-golang/iocli/gen/marker"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &commonCodeGenerationPlugin_{}
		},
	})
	commonCodeGenerationPluginStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &commonCodeGenerationPlugin{}
		},
		ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
			impl := i.(*commonCodeGenerationPlugin)
			var constructFunc commonCodeGenerationPluginConstructFunc = create
			return constructFunc(impl)
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"allimpls": map[string]interface{}{
					"autowireType": "normal",
				},
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(plugin.CodeGeneratorPluginForOneStruct),
					},
				},
			},
		},
	}
	allimpls.RegisterStructDescriptor(commonCodeGenerationPluginStructDescriptor)
	iocGolangAutowireImplmentsAutoInjectionMarkerStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &iocGolangAutowireImplmentsAutoInjectionMarker{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(marker.DefinitionGetter),
					},
				},
			},
		},
		DisableProxy: true,
	}
	allimpls.RegisterStructDescriptor(iocGolangAutowireImplmentsAutoInjectionMarkerStructDescriptor)
	iocGolangAutowireActiveProfileutoInjectionMarkerStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &iocGolangAutowireActiveProfileutoInjectionMarker{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(marker.DefinitionGetter),
					},
				},
			},
		},
		DisableProxy: true,
	}
	allimpls.RegisterStructDescriptor(iocGolangAutowireActiveProfileutoInjectionMarkerStructDescriptor)
	iocGolangAutowireLoadAtOnceMarkerStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &iocGolangAutowireLoadAtOnceMarker{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(marker.DefinitionGetter),
					},
				},
			},
		},
		DisableProxy: true,
	}
	allimpls.RegisterStructDescriptor(iocGolangAutowireLoadAtOnceMarkerStructDescriptor)
}

type commonCodeGenerationPluginConstructFunc func(impl *commonCodeGenerationPlugin) (*commonCodeGenerationPlugin, error)
type commonCodeGenerationPlugin_ struct {
	Name_                           func() string
	Type_                           func() plugin.Type
	Init_                           func(markers markers.MarkerValues)
	GenerateSDMetadataForOneStruct_ func(w plugin.CodeWriter)
	GenerateInFileForOneStruct_     func(w plugin.CodeWriter)
}

func (c *commonCodeGenerationPlugin_) Name() string {
	return c.Name_()
}

func (c *commonCodeGenerationPlugin_) Type() plugin.Type {
	return c.Type_()
}

func (c *commonCodeGenerationPlugin_) Init(markers markers.MarkerValues) {
	c.Init_(markers)
}

func (c *commonCodeGenerationPlugin_) GenerateSDMetadataForOneStruct(w plugin.CodeWriter) {
	c.GenerateSDMetadataForOneStruct_(w)
}

func (c *commonCodeGenerationPlugin_) GenerateInFileForOneStruct(w plugin.CodeWriter) {
	c.GenerateInFileForOneStruct_(w)
}

type commonCodeGenerationPluginIOCInterface interface {
	Name() string
	Type() plugin.Type
	Init(markers markers.MarkerValues)
	GenerateSDMetadataForOneStruct(w plugin.CodeWriter)
	GenerateInFileForOneStruct(w plugin.CodeWriter)
}

var _commonCodeGenerationPluginSDID string
var _iocGolangAutowireImplmentsAutoInjectionMarkerSDID string
var _iocGolangAutowireActiveProfileutoInjectionMarkerSDID string
var _iocGolangAutowireLoadAtOnceMarkerSDID string
