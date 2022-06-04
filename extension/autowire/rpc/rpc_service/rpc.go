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

package rpc_service

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		rpcAutowire := &Autowire{}
		// todo parse rpc param
		rpcAutowire.Autowire = singleton.NewSingletonAutowire(nil, nil, rpcAutowire)
		return rpcAutowire
	}())
}

const Name = "rpc-service"

type Autowire struct {
	autowire.Autowire
}

func (a *Autowire) TagKey() string {
	return Name
}

// GetAllStructDescriptors re-write SingletonAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return rpcStructDescriptorMap
}

var rpcStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	s.SetAutowireType(Name)
	sdID := s.ID()
	s.ParamFactory = func() interface{} {
		return &Param{}
	}
	s.ConstructFunc = func(impl interface{}, p interface{}) (interface{}, error) {
		if p == nil {
			// param not configured in server side, set default param
			p = &Param{
				ExportPort: "2022",
			}
		}
		param := p.(*Param)
		iocProtocolImpl, err := protocol_impl.GetIOCProtocol(&protocol_impl.Param{
			ExportPort: param.ExportPort,
		})
		if err != nil {
			return nil, err
		}
		_, err = protocol_impl.ServiceMap.Register(sdID, protocol_impl.IOCProtocolName, "", "", impl)
		if err != nil {
			panic(err)
		}

		invURL, _ := common.NewURL(protocol_impl.IOCProtocolName+"://",
			common.WithParamsValue(constant.InterfaceKey, sdID),
			common.WithParamsValue("alias", s.Alias),
		)
		defaultProxyInvoker := newProxyInvoker(invURL)
		iocProtocolImpl.Export(defaultProxyInvoker)
		return impl, nil
	}
	rpcStructDescriptorMap[sdID] = s
}

func GetImpl(key string) (interface{}, error) {
	return autowire.Impl(Name, key, nil)
}
