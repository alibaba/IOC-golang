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

package dynamic_plugin

import (
	"context"
	"fmt"
	"os"
	"plugin"

	"github.com/alibaba/ioc-golang/logger"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/extension/aop/dynamic_plugin/api/ioc_golang/aop/dynamic_plugin"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

// todo: rollback support: we should cache dynamic plugins and origin impl to support rollback
type dynamicPluginServiceImpl struct {
	dynamic_plugin.UnimplementedDynamicPluginServiceServer
}

func (d *dynamicPluginServiceImpl) Update(ctx context.Context, req *dynamic_plugin.DynamicPluginUpdateRequest) (*dynamic_plugin.DynamicPluginUpdateResponse, error) {
	tempFile, err := os.CreateTemp("", req.GetPluginName())
	defer func() {
		_ = tempFile.Close()
	}()
	if err != nil {
		return &dynamic_plugin.DynamicPluginUpdateResponse{
			Message: err.Error(),
			Success: false,
		}, err
	}
	if _, err = tempFile.Write(req.PluginFile); err != nil {
		return &dynamic_plugin.DynamicPluginUpdateResponse{
			Message: err.Error(),
			Success: false,
		}, err
	}
	if err := updateDynamicPlugin(req.AutowireType, req.GetSdid(), tempFile.Name(), req.GetPluginName()); err != nil {
		return &dynamic_plugin.DynamicPluginUpdateResponse{
			Message: err.Error(),
			Success: false,
		}, err
	}
	return &dynamic_plugin.DynamicPluginUpdateResponse{
		Success: true,
	}, nil
}

func updateDynamicPlugin(autowireType, sdid string, pluginPath, pluginName string) error {
	// 1. get struct descriptor
	structDescriptor := autowire.GetStructDescriptor(sdid)
	if structDescriptor == nil {
		errMsg := fmt.Sprintf("To update SDID %s struct descriptor not found, update failed", sdid)
		logger.Red("[AOP] Dynamic plugin " + errMsg)
		return fmt.Errorf(errMsg)
	}

	// 2. load plugin
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	pluginImpl, err := plug.Lookup(pluginName)
	if err != nil {
		return err
	}

	// todo check if pluginImpl impls target interface, if not, return with error

	// 3. rewrite struct descriptor
	structDescriptor.Factory = func() interface{} {
		return pluginImpl
	}
	// todo construct function have struct-level operation, so it should not run with plugin struct
	structDescriptor.ConstructFunc = nil
	structDescriptor.SDID = sdid

	// 4. register to target autowire type
	switch autowireType {
	case singleton.Name:
		singleton.RegisterStructDescriptor(structDescriptor)
		// get constructed struct
		constructedPluginImpl, err := autowire.ImplByForce(singleton.Name, sdid, nil)
		if err != nil {
			return err
		}
		// get old proxy wrapper struct
		proxy, err := singleton.GetImplWithProxy(sdid, nil)
		if err != nil {
			return err
		}
		// redirect old proxy wrapper struct to new constructed plugin impl
		if err := autowire.GetProxyImplFunction()(constructedPluginImpl, proxy, sdid); err != nil {
			return err
		}
	case normal.Name:
		normal.RegisterStructDescriptor(structDescriptor)
	default:
		return fmt.Errorf("autowire type %s not support plugin dynamic update", autowireType)
	}
	return nil
}
