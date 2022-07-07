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
	Name                string
	Interceptor         Interceptor
	RPCInterceptor      RPCInterceptor
	GRPCServiceRegister gRPCServiceRegister
	ConfigLoader        func(config *common.Config)
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

type gRPCServiceRegister func(server *grpc.Server)

var aops = make([]AOP, 0)
var interceptors = make([]Interceptor, 0)
var rpcInterceptors = make([]RPCInterceptor, 0)
var grpcServiceRegisters = make([]gRPCServiceRegister, 0)
var configLoaderFuncs = make([]common.ConfigLoader, 0)

func RegisterAOP(aopImpl AOP) {
	aops = append(aops, aopImpl)
	if aopImpl.Interceptor != nil {
		interceptors = append(interceptors, aopImpl.Interceptor)
	}
	if aopImpl.RPCInterceptor != nil {
		rpcInterceptors = append(rpcInterceptors, aopImpl.RPCInterceptor)
	}
	if aopImpl.GRPCServiceRegister != nil {
		grpcServiceRegisters = append(grpcServiceRegisters, aopImpl.GRPCServiceRegister)
	}
	if aopImpl.ConfigLoader != nil {
		configLoaderFuncs = append(configLoaderFuncs, aopImpl.ConfigLoader)
	}
}

func GetRPCInterceptors() []RPCInterceptor {
	return rpcInterceptors
}
