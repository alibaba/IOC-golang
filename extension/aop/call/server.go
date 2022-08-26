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

package call

import (
	"context"
	"encoding/json"
	"fmt"

	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/aop/call/api/ioc_golang/aop/call"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/proxy"
	"github.com/alibaba/ioc-golang/logger"
)

const proxyProtocol = "aop" + Name

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

type callServiceImpl struct {
	call.UnimplementedCallServiceServer
}

func (l *callServiceImpl) Call(_ context.Context, request *call.CallRequest) (*call.CallResponse, error) {
	logger.Red("[Debug Server] Receive call request %+v\n", request.String())
	impl, err := autowire.ImplWithProxy(request.GetAutowireType(), request.GetSdid(), nil)
	if err != nil {
		errMsg := fmt.Sprintf("[AOP call] Get impl with autowire type %s, sdid = %s failed with error = %s", request.GetAutowireType(), request.GetSdid(), err)
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 1. register service metadata with aopCall protocol
	_, _ = proxy.MetadataMap.Register(request.GetSdid(), proxyProtocol, "", "", impl)

	// 2. get registered service metadata with aopCall protocol
	m, ok := proxy.MetadataMap.GetServiceByServiceKey(proxyProtocol, request.GetSdid()).Method()[request.GetMethodName()]
	if !ok {
		errMsg := fmt.Sprintf("[AOP call] Call function with autowire type %s, sdid = %s, method= %s failed, method not found",
			request.GetAutowireType(), request.GetSdid(), request.GetMethodName())
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 3. parse request params if necessary
	arguments := make([]interface{}, 0)
	if rawParams := request.GetParams(); len(rawParams) > 0 {
		arguments, err = protocol_impl.ParseArgs(m.ArgsType(), rawParams)
		if err != nil {
			errMsg := fmt.Sprintf("[AOP call] Call function with autowire type %s, sdid = %s, method= %s failed, parse arguments failed with error = %s",
				request.GetAutowireType(), request.GetSdid(), request.GetMethodName(), err.Error())
			logger.Red(errMsg)
			return nil, fmt.Errorf(errMsg)
		}
	}

	// 4. do invoke
	rsp := proxy.NewProxyInvoker(proxyProtocol, request.GetSdid(), "").Invoke(context.Background(),
		invocation.NewRPCInvocation(request.GetMethodName(), arguments, nil)).Result()

	// 5. marshal response
	rspData, err := json.Marshal(rsp)
	if err != nil {
		errMsg := fmt.Sprintf("[AOP call] Call function with autowire type %s, sdid = %s, method= %s response marshal failed with error = %s",
			request.GetAutowireType(), request.GetSdid(), request.GetMethodName(), err.Error())
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return &call.CallResponse{
		Params:       request.GetParams(),
		ReturnValues: rspData,
		Sdid:         request.GetSdid(),
		MethodName:   request.GetMethodName(),
	}, nil
}
