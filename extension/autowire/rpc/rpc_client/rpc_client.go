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
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		rpcAutowire := &Autowire{}
		// todo parse rpc param
		rpcAutowire.Autowire = normal.NewNormalAutowire(&sdidParser{}, getDefaultParamLoader(), rpcAutowire)
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
	toInvokeSDID, err := util.ToRPCClientStubInterfaceSDID(s.ID())
	if err != nil {
		panic(err)
	}
	s.ParamFactory = func() interface{} {
		return &Param{}
	}
	s.ConstructFunc = func(impl interface{}, p interface{}) (interface{}, error) {
		// client side load as injection, meaningful
		param := p.(*Param)
		iocProtocolInterface, err := protocol_impl.GetIOCProtocolIOCInterface(&protocol_impl.Param{
			Address: param.Address,
			Timeout: param.Timeout,
		})
		if err != nil {
			return nil, err
		}
		newProxy := proxy.NewProxy(newProxyInvoker(iocProtocolInterface, toInvokeSDID), nil, nil)
		defaultProxyImplementFunc(newProxy, impl)
		return impl, nil
	}

	rpcClientStructDescriptorMap[s.ID()] = s
	autowire.RegisterStructDescriptor(s)
}

/*
GetImpl returns impl ptr of rpc client, key should has lowercase character prefix, and 'IOCRPCClient' suffix
like 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient'

The returned interface is proxy struct pointer
like pointer of 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.simpleServiceIOCRPCClient_'

The returned interface can be asserted to ioc rpc client interface
like 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient'

An example to use this API is :
```go
simpleClient, err := rpc_client.GetImpl("github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient", param)
if err != nil{
    panic(err)
}
usr, err := simpleClient.(api.SimpleServiceIOCRPCClient).GetUser("laurence", 23)
```
*/
func GetImpl(key string, param *Param) (interface{}, error) {
	clientStubStructkey, err := util.ToRPCClientStubSDID(key)
	if err != nil {
		return nil, err
	}
	return autowire.ImplWithProxy(Name, clientStubStructkey, param)
}

/*
ImplClientStub returns impl ptr of rpc client, clientStubPtr should be a generated ioc rpc interface client ptr
like 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient'

new(simpleServiceIOCRPCClient)

The returned interface is proxy struct pointer
like pointer of 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.simpleServiceIOCRPCClient_'

The returned interface can be asserted to ioc rpc client interface
like 'github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient'

An example to use this API is :
```go
import(
	"github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/service/api"
)

simpleClient, err := rpc_client.ImplClientStub(new(api.SimpleServiceIOCRPCClient), param)
if err != nil{
    panic(err)
}
usr, err := simpleClient.(api.SimpleServiceIOCRPCClient).GetUser("laurence", 23)
```
*/
func ImplClientStub(clientStubPtr interface{}, param *Param) (interface{}, error) {
	clientStubType := util.GetTypeFromInterface(clientStubPtr)
	return GetImpl(clientStubType.PkgPath()+"."+clientStubType.Name(), param)
}
