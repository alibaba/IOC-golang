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

package proxy

import (
	"context"
	"reflect"
	"strings"

	"github.com/alibaba/ioc-golang/autowire"

	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	perrors "github.com/pkg/errors"
)

func NewProxyInvoker(proto, sdid, alias string) protocol.Invoker {
	invURL, _ := common.NewURL(proto+"://",
		common.WithParamsValue(constant.InterfaceKey, sdid),
		common.WithParamsValue(autowire.AliasKey, alias),
	)
	return &ProxyInvoker{
		BaseInvoker: *protocol.NewBaseInvoker(invURL),
	}
}

// ProxyInvoker is a invoker struct
type ProxyInvoker struct {
	protocol.BaseInvoker
}

// Invoke is used to call service method by invocation
func (pi *ProxyInvoker) Invoke(ctx context.Context, invocation protocol.Invocation) protocol.Result {
	result := &protocol.RPCResult{}
	result.SetAttachments(invocation.Attachments())

	// get providerUrl. The origin url may be is registry URL.
	url := getProviderURL(pi.GetURL())

	methodName := invocation.MethodName()
	proto := url.Protocol
	path := strings.TrimPrefix(url.Path, "/")
	args := invocation.Arguments()

	// get service
	svc := MetadataMap.GetServiceByServiceKey(proto, url.ServiceKey())
	if svc == nil {
		result.SetError(perrors.Errorf("cannot find service [%s] in %s", path, proto))
		return result
	}

	// get method
	method := svc.Method()[methodName]
	if method == nil {
		result.SetError(perrors.Errorf("cannot find method [%s] of service [%s] in %s", methodName, path, proto))
		return result
	}

	in := []reflect.Value{svc.Rcvr()}
	if method.CtxType() != nil {
		ctx = context.WithValue(ctx, constant.AttachmentKey, invocation.Attachments())
		in = append(in, method.SuiteContext(ctx))
	}

	// prepare argv
	if (len(method.ArgsType()) == 1 || len(method.ArgsType()) == 2 && method.ReplyType() == nil) && method.ArgsType()[0].String() == "[]interface {}" {
		in = append(in, reflect.ValueOf(args))
	} else {
		for i := 0; i < len(args); i++ {
			t := reflect.ValueOf(args[i])
			if !t.IsValid() {
				at := method.ArgsType()[i]
				if at.Kind() == reflect.Ptr {
					at = at.Elem()
				}
				t = reflect.New(at)
			}
			in = append(in, t)
		}
	}
	returnValues := method.Method().Func.Call(in)
	returnInterfaces := make([]interface{}, 0)
	for _, v := range returnValues {
		returnInterfaces = append(returnInterfaces, v.Interface())
	}

	var retErr interface{}
	if len(returnValues) >= 1 {
		retErr = returnValues[len(returnValues)-1].Interface()
	}
	if err, ok := retErr.(error); ok && err != nil {
		// error -> string for transfer
		returnInterfaces[len(returnInterfaces)-1] = err.Error()
	}
	result.SetResult(returnInterfaces)
	return result
}

func getProviderURL(url *common.URL) *common.URL {
	if url.SubURL == nil {
		return url
	}
	return url.SubURL
}
