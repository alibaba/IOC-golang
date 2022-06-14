package rpc_client

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/param_loader"
)

var defaultParamLoaderSingleton autowire.ParamLoader

func getDefaultParamLoader() autowire.ParamLoader {
	if defaultParamLoaderSingleton == nil {
		defaultParamLoaderSingleton = &paramLoader{
			defaultConfigParamLoader:     getConfigParamLoader(),
			defaultTagParamLoader:        param_loader.GetDefaultTagParamLoader(),
			defaultTagPointToParamLoader: getTagPointToConfigParamLoader(),
		}
	}
	return defaultParamLoaderSingleton
}

type paramLoader struct {
	defaultConfigParamLoader     autowire.ParamLoader
	defaultTagParamLoader        autowire.ParamLoader
	defaultTagPointToParamLoader autowire.ParamLoader
}

func (d *paramLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if param, err := d.defaultTagPointToParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warning
	if param, err := d.defaultTagParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	// todo log warning

	return d.defaultConfigParamLoader.Load(sd, fi)
}
