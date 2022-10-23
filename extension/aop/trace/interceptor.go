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
	"fmt"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
	traceCommon "github.com/alibaba/ioc-golang/extension/aop/trace/common"
	"github.com/alibaba/ioc-golang/extension/aop/trace/goroutine_trace"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy:autoInjection=false

type traceInterceptor struct {
	GoRoutineInterceptor        goroutine_trace.GoRoutineTraceInterceptorIOCInterface `singleton:""`
	debugServerTraceByMethodCtx *debugServerTraceByMethodContext
}

func (m *traceInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {
	// 1. find if already in goroutine tracing
	if m.GoRoutineInterceptor.GetCurrentGRTracingContext(traceGoRoutineInterceptorFacadeCtxType) != nil {
		m.GoRoutineInterceptor.BeforeInvoke(ctx, traceGoRoutineInterceptorFacadeCtxType)
		return
	}
	// current invocation not in goroutine tracing

	// 2. try to get matched method tracing context
	debugServerTraceByMethodCtx := m.debugServerTraceByMethodCtx
	if debugServerTraceByMethodCtx == nil {
		return
	}
	// method tracing found,
	if debugServerTraceByMethodCtx.fieldMatcher != nil && !debugServerTraceByMethodCtx.fieldMatcher.Match(ctx.Params) {
		// doesn't match trace by method
		return
	}
	if debugServerTraceByMethodCtx.sdid != "" && debugServerTraceByMethodCtx.sdid != ctx.SDID {
		// doesn't match sdid
		return
	}
	if debugServerTraceByMethodCtx.methodName != "" && debugServerTraceByMethodCtx.methodName != ctx.MethodName {
		// doesn't match method
		return
	}
	// match method tracing context found

	// 3.start goroutine tracing
	// create facade ctx
	facadeCtx, err := GettraceGoRoutineInterceptorFacadeCtx(&traceGoRoutineInterceptorFacadeCtxParam{
		trace:     newTrace(ctx.MethodFullName),
		maxDepth:  debugServerTraceByMethodCtx.maxDepth,
		maxLength: debugServerTraceByMethodCtx.maxLength,
	})
	if err != nil {
		log.Printf("traceInterceptor GettraceGoRoutineInterceptorFacadeCtx failed with erorr = %s\n", err.Error())
		return
	}
	// create gr trace ctx
	grCtx, _ := goroutine_trace.GetGoRoutineTracingContext(&goroutine_trace.GoRoutineTracingContextParams{
		FacadeCtx:              facadeCtx,
		EntranceMethodFullName: ctx.MethodFullName,
	})

	// start tracing
	m.GoRoutineInterceptor.AddCurrentGRTracingContext(grCtx)
	m.GoRoutineInterceptor.BeforeInvoke(ctx, traceGoRoutineInterceptorFacadeCtxType)
}

func (m *traceInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	m.GoRoutineInterceptor.AfterInvoke(ctx, traceGoRoutineInterceptorFacadeCtxType)
}

func (m *traceInterceptor) StartTraceByMethod(traceCtx *debugServerTraceByMethodContext) {
	m.debugServerTraceByMethodCtx = traceCtx
}

func (m *traceInterceptor) StopTraceByMethod() {
	m.debugServerTraceByMethodCtx = nil
}

func (m *traceInterceptor) GetCurrentSpan() opentracing.Span {
	if currentGRTracingCtx := m.GoRoutineInterceptor.GetCurrentGRTracingContext(traceGoRoutineInterceptorFacadeCtxType); currentGRTracingCtx != nil {
		facadeCtx := currentGRTracingCtx.GetFacadeCtx().(*traceGoRoutineInterceptorFacadeCtx)
		return facadeCtx.trace.currentSpan.span
	}
	return nil
}

var valueDepth = traceCommon.DefaultRecordValuesDepth
var valueLength = traceCommon.DefaultRecordValuesLength

const traceGoRoutineInterceptorFacadeCtxType = "traceGoRoutineInterceptorFacadeCtx"

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:constructFunc=newTraceGoRoutineInterceptorFacadeCtx
// +ioc:autowire:paramType=traceGoRoutineInterceptorFacadeCtxParam

type traceGoRoutineInterceptorFacadeCtx struct {
	traceGoRoutineInterceptorFacadeCtxParam
}

type traceGoRoutineInterceptorFacadeCtxParam struct {
	trace *trace

	// must be set for rpc interceptor
	clientSpanContext opentracing.SpanContext

	// optional
	maxDepth  int64
	maxLength int64
}

func (p *traceGoRoutineInterceptorFacadeCtxParam) newTraceGoRoutineInterceptorFacadeCtx(c *traceGoRoutineInterceptorFacadeCtx) (*traceGoRoutineInterceptorFacadeCtx, error) {
	if p.trace == nil {
		return nil, fmt.Errorf("traceGoRoutineInterceptorFacadeCtx param traceGoRoutineInterceptorFacadeCtxParam field trace is nil")
	}
	if p.maxDepth == 0 {
		p.maxDepth = traceCommon.DefaultRecordValuesDepth
	}
	if p.maxLength == 0 {
		p.maxLength = traceCommon.DefaultRecordValuesLength
	}
	c.traceGoRoutineInterceptorFacadeCtxParam = *p
	return c, nil
}

func (t *traceGoRoutineInterceptorFacadeCtx) BeforeInvoke(ctx *aop.InvocationContext) {
	currentSpan := t.trace.addChildSpan(ctx.MethodFullName)
	currentSpan.span.LogFields(opentracingLog.String(traceCommon.SpanParamsKey, common.ReflectValues2String(ctx.Params, valueDepth, valueLength)))
}

func (t *traceGoRoutineInterceptorFacadeCtx) AfterInvoke(ctx *aop.InvocationContext) {
	currentSpan := t.trace.currentSpan.span
	t.trace.returnSpan()
	currentSpan.LogFields(opentracingLog.String(traceCommon.SpanReturnValuesKey, common.ReflectValues2String(ctx.ReturnValues, int(t.maxDepth), int(t.maxLength))))
}
func (t *traceGoRoutineInterceptorFacadeCtx) Type() string {
	return traceGoRoutineInterceptorFacadeCtxType
}
