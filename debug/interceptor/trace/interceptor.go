package trace

import (
	"sync"

	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/debug/interceptor"
	"github.com/alibaba/ioc-golang/debug/interceptor/common"
)

type MethodTraceInterceptor struct {
	tracingMethodMap     sync.Map // watch stores methodUniqueKey -> TraceContext
	goRoutineInterceptor *goRoutineTraceInterceptor
}

func (m *MethodTraceInterceptor) BeforeInvoke(ctx *interceptor.InvocationContext) {
	// 1. find if already in goroutine tracing
	if m.goRoutineInterceptor.IsCurrentGRTracing() {
		m.goRoutineInterceptor.BeforeInvoke(ctx)
		return
	}
	// temp invocation not in goroutine tracing

	// 2. try to get matched method tracing context
	watchCtxInterface, ok := m.tracingMethodMap.Load(common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName))
	if !ok {
		// method tracing context not found, no need to trace
		return
	}
	// method tracing found,
	methodTracingCtx := watchCtxInterface.(*methodTracingContext)
	if methodTracingCtx.fieldMatcher != nil && !methodTracingCtx.fieldMatcher.Match(ctx.Params) {
		// doesn't match trace by method
		return
	}
	// match method tracing context found

	// 3.start goroutine tracing
	grCtx := NewGoRoutineTracingContext(common.CurrentCallingMethodName())
	methodTracingCtx.addGoroutineTraceContext(grCtx)
	m.goRoutineInterceptor.AddCurrentGRTracingContext(grCtx)
	m.goRoutineInterceptor.BeforeInvoke(ctx)
}

func (m *MethodTraceInterceptor) AfterInvoke(ctx *interceptor.InvocationContext) {
	m.goRoutineInterceptor.AfterInvoke(ctx)
}

func (m *MethodTraceInterceptor) StartTraceByMethod(traceCtx *methodTracingContext) {
	methodUniqueKey := common.GetMethodUniqueKey(traceCtx.sdid, traceCtx.methodName)
	// FIXME: Now we only support one watcher during whole rpc links
	m.tracingMethodMap.Store(methodUniqueKey, traceCtx)
}

func (m *MethodTraceInterceptor) StopTraceByMethod(traceCtx *methodTracingContext) {
	methodUniqueKey := common.GetMethodUniqueKey(traceCtx.sdid, traceCtx.methodName)
	m.tracingMethodMap.Delete(methodUniqueKey)
}

// TraceCurrentGR is used in rpc-server side, to continue tracing.
func (m *MethodTraceInterceptor) TraceCurrentGR(traceCtx *goRoutineTracingContext) {
	m.goRoutineInterceptor.AddCurrentGRTracingContext(traceCtx)
}

// StopTraceCurrentGR is used in rpc-server side, to continue tracing.
func (m *MethodTraceInterceptor) StopTraceCurrentGR() {
	m.goRoutineInterceptor.DeleteCurrentGRTracingContext()
}

func (m *MethodTraceInterceptor) GetCurrentSpan() opentracing.Span {
	currentGRTracingCtx := m.goRoutineInterceptor.GetCurrentGRTracingContext()
	if currentGRTracingCtx != nil {
		return currentGRTracingCtx.getTrace().currentSpan.span
	}
	return nil
}

var methodTraceInterceptorSingleton *MethodTraceInterceptor

func GetTraceInterceptor() *MethodTraceInterceptor {
	if methodTraceInterceptorSingleton == nil {
		methodTraceInterceptorSingleton = &MethodTraceInterceptor{
			goRoutineInterceptor: GetGoRoutineTraceInterceptor(),
		}
	}
	return methodTraceInterceptorSingleton
}
