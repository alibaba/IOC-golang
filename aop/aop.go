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

package aop

import (
	"net/http"

	"github.com/alibaba/ioc-golang/aop/common"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type AOP struct {
	Name string
	// ConfigLoader is called during ioc.Load() when aop is enabled
	ConfigLoader func(config *common.Config)

	// InterceptorFactory is called after ConfigLoader is called, when bot aop and debug-server are enabled
	InterceptorFactory interceptorFactory
	// RPCInterceptorFactory is called after ConfigLoader is called, when bot aop and debug-server are enabled
	RPCInterceptorFactory rpcInterceptorFactory
	// GRPCServiceRegister is called after ConfigLoader is called, when bot aop and debug-server are enabled
	GRPCServiceRegister gRPCServiceRegister
}

type Interceptor interface {
	BeforeInvoke(ctx *InvocationContext)
	AfterInvoke(ctx *InvocationContext)
}

type RPCInterceptor interface {
	BeforeClientInvoke(req *http.Request) error
	AfterClientInvoke(rsp *http.Response) error
	BeforeServerInvoke(c *gin.Context) error
	AfterServerInvoke(c *gin.Context) error
}

type interceptorFactory func() Interceptor
type rpcInterceptorFactory func() RPCInterceptor
type gRPCServiceRegister func(server *grpc.Server)

var aops = make([]AOP, 0)

var rpcInterceptors []RPCInterceptor
var interceptors []Interceptor

var interceptorFactories = make([]interceptorFactory, 0)
var rpcInterceptorFactories = make([]rpcInterceptorFactory, 0)
var grpcServiceRegisters = make([]gRPCServiceRegister, 0)
var configLoaderFuncs = make([]common.ConfigLoader, 0)

func RegisterAOP(aopImpl AOP) {
	aops = append(aops, aopImpl)
	if aopImpl.InterceptorFactory != nil {
		interceptorFactories = append(interceptorFactories, aopImpl.InterceptorFactory)
	}
	if aopImpl.RPCInterceptorFactory != nil {
		rpcInterceptorFactories = append(rpcInterceptorFactories, aopImpl.RPCInterceptorFactory)
	}
	if aopImpl.GRPCServiceRegister != nil {
		grpcServiceRegisters = append(grpcServiceRegisters, aopImpl.GRPCServiceRegister)
	}
	if aopImpl.ConfigLoader != nil {
		configLoaderFuncs = append(configLoaderFuncs, aopImpl.ConfigLoader)
	}
}

func GetRPCInterceptors() []RPCInterceptor {
	if !enabled {
		return make([]RPCInterceptor, 0)
	}
	if rpcInterceptors == nil {
		rpcInterceptors = make([]RPCInterceptor, 0)
		for _, f := range rpcInterceptorFactories {
			rpcInterceptors = append(rpcInterceptors, f())
		}
	}
	return rpcInterceptors
}

func getInterceptors() []Interceptor {
	if !enabled {
		return make([]Interceptor, 0)
	}
	if interceptors == nil {
		interceptors = make([]Interceptor, 0)
		for _, f := range interceptorFactories {
			interceptors = append(interceptors, f())
		}
	}
	return interceptors
}

var enabled = false

// enableAOP let all aop load config and set enable to true
func enableAOP(aopConfig *common.Config) {
	// let aop load config
	for _, cl := range configLoaderFuncs {
		cl(aopConfig)
	}

	enabled = true
}
