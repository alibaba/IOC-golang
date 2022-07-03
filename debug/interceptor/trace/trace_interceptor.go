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
	"sync"

	"github.com/alibaba/ioc-golang/debug/interceptor"

	"github.com/opentracing/opentracing-go"
	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/debug/interceptor/common"
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

func (t *Interceptor) BeforeInvoke(ctx *interceptor.InvocationContext) {
	// 1. if current goroutine is watched?
	grID := goid.Get()
	if val, ok := t.tracingGrIDMap.Load(grID); ok {
		// this goRoutine is watched, add new child node
		val.(*Context).getTrace(grID).addChildSpan(common.CurrentCallingMethodName())
		return
	}

	// 2. try to get matched watchCtx
	watchCtxInterface, ok := t.tracingMethodMap.Load(common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName))
	if !ok {
		return
	}
	watchCtx := watchCtxInterface.(*Context)
	if watchCtx.FieldMatcher != nil && !watchCtx.FieldMatcher.Match(ctx.Params) {
		// doesn't match
		return
	}
	if watchCtx.targetGR && grID != watchCtx.grID {
		return
	}

	// 3.start gr tracing
	watchCtx.createTrace(grID, common.CurrentCallingMethodName())
	t.tracingGrIDMap.Store(grID, watchCtx)
}

func (t *Interceptor) AfterInvoke(_ *interceptor.InvocationContext) {
	// if current goroutine is watched?
	grID := goid.Get()
	if val, ok := t.tracingGrIDMap.Load(grID); ok {
		// this goRoutine is watched, return span
		traceCtx := val.(*Context)
		traceCtx.getTrace(grID).returnSpan()

		// calculate level
		if common.TraceLevel(traceCtx.getTrace(grID).entranceMethod) == 0 {
			traceCtx.finish(grID)
			t.tracingGrIDMap.Delete(grID)
		}
	}
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

var traceInterceptorSingleton *Interceptor

func GetTraceInterceptor() *Interceptor {
	if traceInterceptorSingleton == nil {
		traceInterceptorSingleton = &Interceptor{}
	}
	return traceInterceptorSingleton
}
