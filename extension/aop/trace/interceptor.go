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
	"github.com/opentracing/opentracing-go"

	traceCommon "github.com/alibaba/ioc-golang/extension/aop/trace/common"

	"github.com/alibaba/ioc-golang/aop"
)

type methodTraceInterceptor struct {
	goRoutineInterceptor *goRoutineTraceInterceptor
	tracingCtx           *methodTracingContext
}

func (m *methodTraceInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {
	// 1. find if already in goroutine tracing
	if m.goRoutineInterceptor.IsCurrentGRTracing() {
		m.goRoutineInterceptor.BeforeInvoke(ctx)
		return
	}
	// current invocation not in goroutine tracing

	// 2. try to get matched method tracing context
	traceCtx := m.tracingCtx
	if traceCtx == nil {
		return
	}
	// method tracing found,
	if traceCtx.fieldMatcher != nil && !traceCtx.fieldMatcher.Match(ctx.Params) {
		// doesn't match trace by method
		return
	}
	if traceCtx.sdid != "" && traceCtx.sdid != ctx.SDID {
		// doesn't match sdid
		return
	}
	if traceCtx.methodName != "" && traceCtx.methodName != ctx.MethodName {
		// doesn't match method
		return
	}
	// match method tracing context found

	// 3.start goroutine tracing
	grCtx := newGoRoutineTracingContext(ctx.MethodFullName)
	m.goRoutineInterceptor.AddCurrentGRTracingContext(grCtx)
	m.goRoutineInterceptor.BeforeInvoke(ctx)
}

func (m *methodTraceInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	m.goRoutineInterceptor.AfterInvoke(ctx)
}

func (m *methodTraceInterceptor) StartTraceByMethod(traceCtx *methodTracingContext) {
	m.tracingCtx = traceCtx
}

func (m *methodTraceInterceptor) StopTraceByMethod() {
	m.tracingCtx = nil
}

// TraceCurrentGR is used in rpc-server side, to continue tracing.
func (m *methodTraceInterceptor) TraceCurrentGR(traceCtx *goRoutineTracingContext) {
	m.goRoutineInterceptor.AddCurrentGRTracingContext(traceCtx)
}

// StopTraceCurrentGR is used in rpc-server side, to continue tracing.
func (m *methodTraceInterceptor) StopTraceCurrentGR() {
	m.goRoutineInterceptor.DeleteCurrentGRTracingContext()
}

func (m *methodTraceInterceptor) GetCurrentSpan() opentracing.Span {
	currentGRTracingCtx := m.goRoutineInterceptor.GetCurrentGRTracingContext()
	if currentGRTracingCtx != nil {
		return currentGRTracingCtx.getTrace().currentSpan.span
	}
	return nil
}

var methodTraceInterceptorSingleton *methodTraceInterceptor

var valueDepth = traceCommon.DefaultRecordValuesDepth

func getTraceInterceptorSingleton() *methodTraceInterceptor {
	if methodTraceInterceptorSingleton == nil {
		methodTraceInterceptorSingleton = &methodTraceInterceptor{
			goRoutineInterceptor: getGoRoutineTraceInterceptor(),
		}
	}
	return methodTraceInterceptorSingleton
}
