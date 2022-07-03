//go:build iocMonkeydebug

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
	"reflect"

	"github.com/alibaba/ioc-golang/debug/interceptor"

	"github.com/glory-go/monkey"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/debug/common"
)

func init() {
	autowire.RegisterMonkeyFunction(implMonkey)
}

// nolint
func implMonkey(servicePtr interface{}, tempInterfaceId string) {
	if util.IsProxyStructPtr(servicePtr) {
		return
	}
	if _, ok := debugMetadata[tempInterfaceId]; !ok {
		debugMetadata[tempInterfaceId] = &common.StructMetadata{
			MethodMetadata: map[string]*common.MethodMetadata{},
		}
	}
	valueOf := reflect.ValueOf(servicePtr)
	typeOf := reflect.TypeOf(servicePtr)
	valueOfElem := valueOf.Elem()
	typeOfElem := valueOfElem.Type()
	if typeOfElem.Kind() != reflect.Struct {
		// Fixme: handle the error gracefully
		panic("invalid struct ptr")
	}

	numField := valueOf.NumMethod()
	for i := 0; i < numField; i++ {
		methodType := typeOf.Method(i)
		if _, ok := debugMetadata[tempInterfaceId].MethodMetadata[methodType.Name]; !ok {
			debugMetadata[tempInterfaceId].MethodMetadata[methodType.Name] = &common.MethodMetadata{}
		}
		if debugMetadata[tempInterfaceId].MethodMetadata[methodType.Name].Guard == nil {
			// each method of one type should only injected once
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(servicePtr), methodType.Name,
				reflect.MakeFunc(methodType.Type, makeCallProxy(tempInterfaceId, methodType.Name, methodType.Type.IsVariadic())).Interface(),
			)
			debugMetadata[tempInterfaceId].MethodMetadata[methodType.Name].Guard = guard
		}
		continue
	}
}

// nolint
func makeCallProxy(sdid, methodName string, isVariadic bool) func(in []reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		debugMetadata[sdid].MethodMetadata[methodName].Lock.Lock()
		debugMetadata[sdid].MethodMetadata[methodName].Guard.Unpatch()
		defer func() {
			debugMetadata[sdid].MethodMetadata[methodName].Guard.Restore()
			debugMetadata[sdid].MethodMetadata[methodName].Lock.Unlock()
		}()
		// interceptor

		invocationCtx := &interceptor.InvocationContext{}
		for _, i := range interceptors {
			if len(in) > 1 {
				invocationCtx = interceptor.NewInvocationContext(nil, sdid, methodName, in[1:])
				i.BeforeInvoke(invocationCtx)
				for idx, p := range invocationCtx.Params {
					in[idx+1] = p
				}
			} else {
				invocationCtx = interceptor.NewInvocationContext(nil, sdid, methodName, []reflect.Value{})
				i.BeforeInvoke(invocationCtx)
			}
		}

		if isVariadic {
			varParam := in[len(in)-1]
			in = in[:len(in)-1]
			for j, l := 0, varParam.Len(); j < l; j++ {
				in = append(in, varParam.Index(j))
			}
		}

		out := make([]reflect.Value, 0)
		if len(in) > 1 {
			out = in[0].MethodByName(methodName).Call(in[1:])
		} else {
			out = in[0].MethodByName(methodName).Call([]reflect.Value{})
		}

		invocationCtx.SetReturnValues(out)
		for _, i := range interceptors {
			i.AfterInvoke(invocationCtx)
		}
		return out
	}
}
