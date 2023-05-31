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

package protocol_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alibaba/ioc-golang/extension/autowire/rpc/proxy"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/logger"

	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"github.com/gin-gonic/gin"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	dubboProtocol "dubbo.apache.org/dubbo-go/v3/protocol"

	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init

// IOCProtocol is ioc protocol impl
type IOCProtocol struct {
	address    string
	exportPort string
	timeout    string
}

func (i *IOCProtocol) Invoke(invocation dubboProtocol.Invocation) dubboProtocol.Result {
	sdID, _ := invocation.GetAttachment("sdid")
	data, _ := json.Marshal(invocation.Arguments())
	invokeURL := DefaultSchema + "://" + i.address + "/" + sdID + "/" + invocation.MethodName()

	timeoutDuration, err := time.ParseDuration(i.timeout)
	if err != nil {
		timeoutDuration = time.Second * 3
	}
	requestCtx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, invokeURL, bytes.NewBuffer(data))
	if err != nil {
		return &dubboProtocol.RPCResult{
			Err: err,
		}
	}

	allRPCInterceptors := aop.GetRPCInterceptors()
	for _, rpcInterceptor := range allRPCInterceptors {
		if err := rpcInterceptor.BeforeClientInvoke(req); err != nil {
			return &dubboProtocol.RPCResult{
				Err: err,
			}
		}
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Red("[IOC Protocol] Invoke %s with error = %s", invokeURL, err)
		return &dubboProtocol.RPCResult{
			Err: err,
		}
	}

	for _, rpcInterceptor := range allRPCInterceptors {
		if err := rpcInterceptor.AfterClientInvoke(rsp); err != nil {
			return &dubboProtocol.RPCResult{
				Err: err,
			}
		}
	}

	rspData, _ := io.ReadAll(rsp.Body)
	replyList := invocation.Reply().(*[]interface{})
	finalIsError := false
	finalErrorNotNil := false
	if length := len(*replyList); length > 0 {
		_, ok := (*replyList)[length-1].(*error)
		if ok {
			finalIsError = true
		}
	}
	err = json.Unmarshal(rspData, replyList)
	if err != nil && finalIsError {
		// error message must be returned
		finalErrorNotNil = true

		// calculate error message detail, try to recover unmarshal failed caused by error not empty, first try to unmarshal to string
		(*replyList)[len(*replyList)-1] = ""
		err = json.Unmarshal(rspData, replyList)
		if err != nil {
			// error is not nil, means previous unmarshal failed because of invalid response, write error message
			err = fmt.Errorf("[IOC Protocol] Unmarshal response from %s with error %s, response data details is %s", invokeURL, err, string(rspData))
			(*replyList)[len(*replyList)-1] = err
		}
		// error is nil means final return value error is returned from server side, and the response is valid
	}
	if err != nil {
		return &dubboProtocol.RPCResult{
			Err: err,
		}
	}
	if finalErrorNotNil {
		realErr := fmt.Errorf((*replyList)[len(*replyList)-1].(string))
		(*replyList)[len(*replyList)-1] = &realErr
	}
	return nil
}

func (i *IOCProtocol) Export(invoker dubboProtocol.Invoker) dubboProtocol.Exporter {
	httpServer := getSingletonGinEngion(i.exportPort)

	sdid := invoker.GetURL().GetParam(constant.InterfaceKey, "")
	clientStubFullName := invoker.GetURL().GetParam(autowire.AliasKey, "")
	svc := proxy.MetadataMap.GetServiceByServiceKey(IOCProtocolName, sdid)
	if svc == nil {
		return nil
	}

	for methodName, methodType := range svc.Method() {
		argsType := methodType.ArgsType()
		tempMethod := methodName
		httpServer.POST(fmt.Sprintf("/%s/%s", clientStubFullName, tempMethod), func(c *gin.Context) {
			allRPCInterceptors := aop.GetRPCInterceptors()
			for _, rpcInterceptor := range allRPCInterceptors {
				if err := rpcInterceptor.BeforeServerInvoke(c); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				}
			}

			reqData, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}
			arguments, err := ParseArgs(argsType, reqData)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}
			rsp := invoker.Invoke(context.Background(),
				invocation.NewRPCInvocation(tempMethod, arguments, nil)).Result()
			c.PureJSON(http.StatusOK, rsp)

			for _, rpcInterceptor := range allRPCInterceptors {
				if err := rpcInterceptor.AfterServerInvoke(c); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				}
			}
		})
	}

	return dubboProtocol.NewBaseExporter(sdid, invoker, nil)
}

func getSingletonGinEngion(exportPort string) *gin.Engine {
	if ginEngionSingleton == nil {
		ginEngionSingleton = gin.Default()
		go func() {
			if err := ginEngionSingleton.Run(":" + exportPort); err != nil {
				// FIXME, should throw error gracefully
				panic(err)
			}
		}()
	}
	return ginEngionSingleton
}

var ginEngionSingleton *gin.Engine

var _ protocol.Protocol = &IOCProtocol{}
