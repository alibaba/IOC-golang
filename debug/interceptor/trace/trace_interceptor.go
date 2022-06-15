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

package trace

import (
	"reflect"
	"runtime"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/debug/interceptor/common"
)

const (
	proxyMethod = "github.com/alibaba/ioc-golang/debug.makeProxyFunction.func1"
)

type Interceptor struct {
	tracingMethodMap sync.Map // watch stores methodUniqueKey -> TraceContext
	tracingGrIDMap   sync.Map // tracingGrIDMap stores goroutine-id -> TraceContext
}

func (t *Interceptor) GetCurrentSpan() opentracing.Span {
	grID := goid.Get()
	val, ok := t.tracingGrIDMap.Load(grID)
	if !ok {
		return nil
	}
	return val.(*Context).getTrace(grID).currentSpan.span
}

func (t *Interceptor) BeforeInvoke(sdid, methodName string, values []reflect.Value) []reflect.Value {
	// 1. if current goroutine is watched?
	grID := goid.Get()
	if val, ok := t.tracingGrIDMap.Load(grID); ok {
		// this goRoutine is watched, add new child node
		val.(*Context).getTrace(grID).addChildSpan(currentCallingMethodName())
		return values
	}

	// 2. try to get matched watchCtx
	watchCtxInterface, ok := t.tracingMethodMap.Load(common.GetMethodUniqueKey(sdid, methodName))
	if !ok {
		return values
	}
	watchCtx := watchCtxInterface.(*Context)
	if watchCtx.FieldMatcher != nil && !watchCtx.FieldMatcher.Match(values) {
		// doesn't match
		return values
	}
	if watchCtx.targetGR && grID != watchCtx.grID {
		// not target gr
		return values
	}

	// 3.start gr tracing
	watchCtx.createTrace(grID, currentCallingMethodName())
	t.tracingGrIDMap.Store(grID, watchCtx)
	return values
}

func (t *Interceptor) AfterInvoke(_, _ string, values []reflect.Value) []reflect.Value {
	// if current goroutine is watched?
	grID := goid.Get()
	if val, ok := t.tracingGrIDMap.Load(grID); ok {
		// this goRoutine is watched, return span
		traceCtx := val.(*Context)
		traceCtx.getTrace(grID).returnSpan()

		// calculate level
		if traceLevel(traceCtx.getTrace(grID).entranceMethod) == 0 {
			traceCtx.finish(grID)
			t.tracingGrIDMap.Delete(grID)
		}
		return values
	}
	return values
}

func (t *Interceptor) Trace(traceCtx *Context) {
	methodUniqueKey := common.GetMethodUniqueKey(traceCtx.SDID, traceCtx.MethodName)
	// FIXME: Now we only support one watcher during whole rpc links
	t.tracingMethodMap.Store(methodUniqueKey, traceCtx)
}

// TraceThisGR is used in rpc-server side, to continue tracing.
func (t *Interceptor) TraceThisGR(traceCtx *Context) {
	traceCtx.targetGR = true
	traceCtx.grID = goid.Get()
	methodUniqueKey := common.GetMethodUniqueKey(traceCtx.SDID, traceCtx.MethodName)
	// FIXME: Now we only support one watcher during whole rpc links
	t.tracingMethodMap.Store(methodUniqueKey, traceCtx)
}

func (t *Interceptor) UnTrace(traceCtx *Context) {
	methodUniqueKey := common.GetMethodUniqueKey(traceCtx.SDID, traceCtx.MethodName)
	t.tracingMethodMap.Delete(methodUniqueKey)
}

func currentCallingMethodName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(4, pc)
	return runtime.FuncForPC(pc[0]).Name()
}

func traceLevel(entranceName string) int64 {
	pc := make([]uintptr, 100)
	n := runtime.Callers(0, pc)
	foundEntrance := false
	level := int64(0)

	for i := n - 1; i >= 0; i-- {
		fName := runtime.FuncForPC(pc[i]).Name()
		if foundEntrance {
			if fName == proxyMethod {
				level++
			}
			continue
		}
		if fName == entranceName {
			foundEntrance = true
		}
	}

	return level - 1
}

var traceInterceptorSingleton *Interceptor

func GetTraceInterceptor() *Interceptor {
	if traceInterceptorSingleton == nil {
		traceInterceptorSingleton = &Interceptor{}
	}
	return traceInterceptorSingleton
}
