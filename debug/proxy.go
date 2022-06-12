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

package debug

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/debug/common"
)

func init() {
	autowire.RegisterProxyFunction(proxyFunction)
}

func proxyFunction(rawPtr interface{}) (interface{}, error) {
	sdid := util.GetSDIDByStructPtr(rawPtr)
	proxySDID := util.GetProxySDIDByStructPtr(rawPtr)
	proxyStructPtr, err := normal.GetImpl(proxySDID, nil)
	if err != nil {
		return nil, err
	}

	if err := implProxy(rawPtr, proxyStructPtr, sdid); err != nil {
		return nil, err
	}
	return proxyStructPtr, nil
}

func implProxy(rawServicePtr, proxyPtr interface{}, sdid string) error {
	if _, ok := debugMetadata[sdid]; !ok {
		debugMetadata[sdid] = &common.StructMetadata{
			MethodMetadata: map[string]*common.MethodMetadata{},
		}
	}
	valueOf := reflect.ValueOf(proxyPtr)
	valueOfElem := valueOf.Elem()
	typeOfElem := valueOfElem.Type()
	if typeOfElem.Kind() != reflect.Struct {
		return fmt.Errorf("invalid struct ptr %+v", proxyPtr)
	}

	valueOfRaw := reflect.ValueOf(rawServicePtr)

	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		methodType := typeOfElem.Field(i)
		f := valueOfElem.Field(i)
		rawMethodName := strings.TrimSuffix(methodType.Name, "_")
		funcRaw := valueOfRaw.MethodByName(rawMethodName)
		// each method of one type should only injected once
		if f.Kind() == reflect.Func && f.IsValid() && f.CanSet() {
			debugMetadata[sdid].MethodMetadata[rawMethodName] = &common.MethodMetadata{}
			f.Set(reflect.MakeFunc(methodType.Type, makeProxyFunction(funcRaw, sdid, rawMethodName, methodType.Type.IsVariadic())))
		}
	}
	return nil
}

func makeProxyFunction(rf reflect.Value, tempInterfaceId, methodName string, isVariadic bool) func(in []reflect.Value) []reflect.Value {
	rawFunction := rf
	return func(in []reflect.Value) []reflect.Value {
		// interceptor
		for _, i := range paramInterceptors {
			in = i.Invoke(tempInterfaceId, methodName, true, in)
		}

		if isVariadic {
			varParam := in[len(in)-1]
			in = in[:len(in)-1]
			for j, l := 0, varParam.Len(); j < l; j++ {
				in = append(in, varParam.Index(j))
			}
		}

		out := rawFunction.Call(in)
		for _, i := range responseInterceptors {
			out = i.Invoke(tempInterfaceId, methodName, false, out)
		}
		return out
	}
}
