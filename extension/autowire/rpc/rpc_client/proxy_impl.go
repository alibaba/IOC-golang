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
	"context"
	"errors"
	"reflect"

	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/proxy"
	invocation_impl "dubbo.apache.org/dubbo-go/v3/protocol/invocation"
)

var typError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()).Type()

// defaultProxyImplementFunc the default function for proxy impl
func defaultProxyImplementFunc(p *proxy.Proxy, v common.RPCService) {
	// check parameters, incoming interface must be a elem's pointer.
	valueOf := reflect.ValueOf(v)

	valueOfElem := valueOf.Elem()

	makeDubboCallProxy := func(methodName string, outs []reflect.Type) func(in []reflect.Value) []reflect.Value {
		return func(in []reflect.Value) []reflect.Value {
			var (
				inv    *invocation_impl.RPCInvocation
				inIArr []interface{}
				inVArr []reflect.Value
			)
			if methodName == "Echo" {
				methodName = "$echo"
			}

			replyInterface := make([]interface{}, 0)
			outReflectValues := make([]reflect.Value, 0)
			for _, o := range outs {
				var reflectValue reflect.Value
				if o.Kind() == reflect.Ptr {
					reflectValue = reflect.New(o.Elem())
				} else {
					reflectValue = reflect.New(o)
				}
				outReflectValues = append(outReflectValues, reflectValue)
				replyInterface = append(replyInterface, reflectValue.Interface())
			}

			start := 0
			end := len(in)
			invCtx := context.Background()
			// retrieve the context from the first argument if existed
			if end > 0 {
				if in[0].Type().String() == "context.Context" {
					if !in[0].IsNil() {
						// the user declared context as method's parameter
						invCtx = in[0].Interface().(context.Context)
					}
					start += 1
				}
			}

			if end-start <= 0 {
				inIArr = []interface{}{}
				inVArr = []reflect.Value{}
			} else if v, ok := in[start].Interface().([]interface{}); ok && end-start == 1 {
				inIArr = v
				inVArr = []reflect.Value{in[start]}
			} else {
				inIArr = make([]interface{}, end-start)
				inVArr = make([]reflect.Value, end-start)
				index := 0
				for i := start; i < end; i++ {
					inIArr[index] = in[i].Interface()
					inVArr[index] = in[i]
					index++
				}
			}

			inv = invocation_impl.NewRPCInvocationWithOptions(invocation_impl.WithMethodName(methodName),
				invocation_impl.WithArguments(inIArr), invocation_impl.WithParameterValues(inVArr))
			inv.SetReply(&replyInterface)

			rpcResult := p.GetInvoker().Invoke(invCtx, inv)

			// interface -> reflectionValues
			reflectValues := make([]reflect.Value, 0)
			returnReflectValues := make([]reflect.Value, 0)
			for idx, v := range replyInterface {
				if v == nil {
					reflectValues = append(reflectValues, outReflectValues[idx])
				} else {
					reflectValues = append(reflectValues, reflect.ValueOf(v))
				}
			}
			for idx, reflectValue := range reflectValues {
				if outs[idx].Kind() != reflect.Ptr {
					returnReflectValues = append(returnReflectValues, reflectValue.Elem())
				} else {
					returnReflectValues = append(returnReflectValues, reflectValue)
				}
			}

			if rpcResult != nil && rpcResult.Error() != nil && len(outs) > 0 && outs[len(outs)-1] == typError {
				returnReflectValues[len(returnReflectValues)-1] = reflect.ValueOf(rpcResult.Error())
			}

			return returnReflectValues
		}
	}

	if err := reflectAndMakeObjectFunc(valueOfElem, makeDubboCallProxy); err != nil {
		return
	}
}

func reflectAndMakeObjectFunc(valueOfElem reflect.Value, makeDubboCallProxy func(methodName string, outs []reflect.Type) func(in []reflect.Value) []reflect.Value) error {
	typeOf := valueOfElem.Type()
	// check incoming interface, incoming interface's elem must be a struct.
	if typeOf.Kind() != reflect.Struct {
		return errors.New("invalid type kind")
	}
	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		t := typeOf.Field(i)
		methodName := t.Name
		f := valueOfElem.Field(i)
		if f.Kind() == reflect.Func && f.IsValid() && f.CanSet() {
			outNum := t.Type.NumOut()

			funcOuts := make([]reflect.Type, outNum)
			for i := 0; i < outNum; i++ {
				funcOuts[i] = t.Type.Out(i)
			}

			// do method proxy here:
			f.Set(reflect.MakeFunc(f.Type(), makeDubboCallProxy(methodName, funcOuts)))
		} else if f.IsValid() && f.CanSet() {
			// for struct combination
			valueOfSub := reflect.New(t.Type)
			valueOfElemInterface := valueOfSub.Elem()
			if valueOfElemInterface.Type().Kind() == reflect.Struct {
				if err := reflectAndMakeObjectFunc(valueOfElemInterface, makeDubboCallProxy); err != nil {
					return err
				}
				f.Set(valueOfElemInterface)
			}
		}
	}
	return nil
}
