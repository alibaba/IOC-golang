package trace

import (
	"sync"

	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
)

type goRoutineTraceInterceptor struct {
	tracingGrIDMap sync.Map // tracingGrIDMap stores goroutine-id -> goRoutineTracingContext
}

func (g *goRoutineTraceInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {
	// 1. if current goroutine is watched?
	if val, ok := g.tracingGrIDMap.Load(ctx.GrID); ok {
		// this goRoutine is watched, add new child node
		val.(*goRoutineTracingContext).getTrace().addChildSpan(common.CurrentCallingMethodName())
		return
	}
}

func (g *goRoutineTraceInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	// if current goroutine is watched?
	if val, ok := g.tracingGrIDMap.Load(ctx.GrID); ok {
		// this goRoutine is watched, return span
		traceCtx := val.(*goRoutineTracingContext)
		traceCtx.getTrace().returnSpan()

		// calculate level
		if common.TraceLevel(traceCtx.getTrace().entranceMethod) == 0 {
			g.tracingGrIDMap.Delete(ctx.GrID)
		}
	}
}

func (g *goRoutineTraceInterceptor) AddCurrentGRTracingContext(ctx *goRoutineTracingContext) {
	grID := goid.Get()
	g.tracingGrIDMap.Store(grID, ctx)
}

func (g *goRoutineTraceInterceptor) DeleteCurrentGRTracingContext() {
	grID := goid.Get()
	g.tracingGrIDMap.Delete(grID)
}

func (g *goRoutineTraceInterceptor) IsCurrentGRTracing() bool {
	grID := goid.Get()
	_, ok := g.tracingGrIDMap.Load(grID)
	return ok
}

func (g *goRoutineTraceInterceptor) GetCurrentGRTracingContext() *goRoutineTracingContext {
	grID := goid.Get()
	val, ok := g.tracingGrIDMap.Load(grID)
	if ok {
		return val.(*goRoutineTracingContext)
	}
	return nil
}

var goRoutineTraceInterceptorSingleton *goRoutineTraceInterceptor

func getGoRoutineTraceInterceptor() *goRoutineTraceInterceptor {
	if goRoutineTraceInterceptorSingleton == nil {
		goRoutineTraceInterceptorSingleton = &goRoutineTraceInterceptor{}
	}
	return goRoutineTraceInterceptorSingleton
}
