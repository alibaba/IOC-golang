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

package rpc_client

import (
	"dubbo.apache.org/dubbo-go/v3/common/proxy"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		rpcAutowire := &Autowire{}
		// todo parse rpc param
		rpcAutowire.Autowire = normal.NewNormalAutowire(nil, nil, rpcAutowire)
		return rpcAutowire
	}())
}

const Name = "rpc-client"

type Autowire struct {
	autowire.Autowire
}

func (a *Autowire) TagKey() string {
	return Name
}

// GetAllStructDescriptors re-write SingletonAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return rpcClientStructDescriptorMap
}

var rpcClientStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	s.SetAutowireType(Name)
	iocRPCClientSDID := s.ID()
	s.ParamFactory = func() interface{} {
		return &Param{}
	}
	s.ConstructFunc = func(impl interface{}, p interface{}) (interface{}, error) {
		// client side load as injection, meaningful
		param := p.(*Param)
		iocProtocolImpl, err := protocol_impl.GetIOCProtocol(&protocol_impl.Param{
			Address: param.Address,
		})
		if err != nil {
			return nil, err
		}
		newProxy := proxy.NewProxy(newProxyInvoker(iocProtocolImpl, iocRPCClientSDID), nil, nil)
		defaultProxyImplementFunc(newProxy, impl)
		return impl, nil
	}

	rpcClientStructDescriptorMap[iocRPCClientSDID] = s
}

func GetImpl(key string) (interface{}, error) {
	return autowire.Impl(Name, key, nil)
}
