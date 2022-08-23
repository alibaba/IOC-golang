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
	"fmt"
	"reflect"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/aop/call/api/ioc_golang/aop/call"
	"github.com/alibaba/ioc-golang/logger"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

type callServiceImpl struct {
	call.UnimplementedCallServiceServer
	allInterfaceMetadataMap common.AllInterfaceMetadata
}

func (l *callServiceImpl) Call(_ context.Context, request *call.CallRequest) (*call.CallResponse, error) {
	logger.Red("[Debug Server] Receive call request %+v\n", request.String())
	impl, err := autowire.ImplWithProxy(request.GetAutowireType(), request.GetSdid(), nil)
	if err != nil {
		errMsg := fmt.Sprintf("[AOP call] Get impl with autowire type %s, sdid = %s failed with error = %s", request.GetAutowireType(), request.GetSdid(), err)
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	valueOfImpl := reflect.ValueOf(impl)
	valueOfElem := valueOfImpl.Elem()

	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		f := valueOfElem.Field(i)
		funcRaw := valueOfImpl.MethodByName(request.MethodName)
		if !(f.Kind() == reflect.Func && f.IsValid()) {
			// not current field
			continue
		}
		// current function is matched, try to call function
		// todo []string -> []reflect.Value

		funcRaw.Call(nil)
	}
	errMsg := fmt.Sprintf("[AOP call] Call function with autowire type %s, sdid = %s, method= %s failed, method not found",
		request.GetAutowireType(), request.GetSdid(), request.GetMethodName())
	logger.Red(errMsg)
	return nil, fmt.Errorf(errMsg)
}
