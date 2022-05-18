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

package interceptor

import (
	"reflect"
	"strings"
	"sync"

	"github.com/alibaba/ioc-golang/debug/common"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

type EditInterceptor struct {
	watchEdit sync.Map
}

func (w *EditInterceptor) Invoke(ctx *common.InterceptorContext, values []reflect.Value) ([]reflect.Value, error) {
	interfaceImplId := ctx.GetSDID()
	methodName := ctx.GetMethod()
	isParam := ctx.IsParam()
	methodUniqueKey := getMethodUniqueKey(ctx)
	watchEditCtxInterface, ok := w.watchEdit.Load(methodUniqueKey)
	if !ok {
		return values, nil
	}
	watchEditCtx := watchEditCtxInterface.(*EditContext)
	if watchEditCtx.FieldMatcher != nil && !watchEditCtx.FieldMatcher.Match(values) {
		// doesn't match
		return values, nil
	}

	// send condition
	sendValues(interfaceImplId, methodName, isParam, values, watchEditCtx.SendCh)

	// block and wait edit signal
	recvMsg := <-watchEditCtx.RecvCh

	// edit
	afterEditedValues, ok := recvMsg.Edit(values)
	if !ok {
		return values, nil
	}
	return afterEditedValues, nil
}

func (w *EditInterceptor) Name() string {
	return "default_edit"
}

type EditContext struct {
	SendCh       chan *boot.WatchResponse
	RecvCh       chan *EditData
	FieldMatcher *FieldMatcher
}

type EditData struct {
	FieldIndex int
	FieldPath  string // A.B.C
	Value      string
}

func (e *EditData) Edit(values []reflect.Value) ([]reflect.Value, bool) {
	if len(values) < e.FieldIndex {
		return nil, false
	}
	targetVal := values[e.FieldIndex]
	valueOfElem := targetVal
	if valueOfElem.Kind() == reflect.Ptr || valueOfElem.Kind() == reflect.Interface {
		valueOfElem = valueOfElem.Elem()
	}
	//typeOfElem := valueOfElem.Type()
	splitedPaths := strings.Split(e.FieldPath, ".")

	for i, p := range splitedPaths {
		val := valueOfElem.FieldByName(p)
		if i == len(splitedPaths)-1 {
			if !val.CanSet() {
				return nil, false
			}
			val.Set(reflect.ValueOf(e.Value))
			return values, true
		}
		if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
			valueOfElem = val.Elem()
		} else {
			valueOfElem = val
		}
	}
	return values, true
}

func (w *EditInterceptor) WatchEdit(ctx *common.InterceptorContext, editCtx *EditContext) {
	methodUniqueKey := getMethodUniqueKey(ctx)
	w.watchEdit.Store(methodUniqueKey, editCtx)
}

func (w *EditInterceptor) UnWatchEdit(ctx *common.InterceptorContext) {
	methodUniqueKey := getMethodUniqueKey(ctx)
	w.watchEdit.Delete(methodUniqueKey)
}

var editInterceptorSingleton *EditInterceptor

func GetEditInterceptor() *EditInterceptor {
	if editInterceptorSingleton == nil {
		editInterceptorSingleton = &EditInterceptor{}
	}
	return editInterceptorSingleton
}
