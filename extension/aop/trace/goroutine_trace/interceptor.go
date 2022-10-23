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

package goroutine_trace

import (
	"sync"

	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy:autoInjection=false

type GoRoutineTraceInterceptor struct {
	tracingGrIDMap sync.Map // tracingGrIDMap stores goroutine-id -> goRoutineTracingContext
}

func (g *GoRoutineTraceInterceptor) BeforeInvoke(ctx *aop.InvocationContext, facadeCtxType string) {
	// if current goroutine is in tracing
	if traceCtxVal, ok := g.tracingGrIDMap.Load(ctx.GrID); ok {
		traceCtx := traceCtxVal.(*GoRoutineTracingContext)
		if traceCtx.GetFacadeCtx().Type() != facadeCtxType {
			return
		}
		traceCtxVal.(*GoRoutineTracingContext).facadeCtx.BeforeInvoke(ctx)
	}
}

func (g *GoRoutineTraceInterceptor) AfterInvoke(ctx *aop.InvocationContext, facadeCtxType string) {
	// if current goroutine is watched?
	if traceCtxVal, ok := g.tracingGrIDMap.Load(ctx.GrID); ok {
		// this goRoutine is in tracing, return span
		traceCtx := traceCtxVal.(*GoRoutineTracingContext)
		if traceCtx.GetFacadeCtx().Type() != facadeCtxType {
			return
		}
		traceCtx.facadeCtx.AfterInvoke(ctx)
		// calculate level
		if common.IsTraceEntrance(traceCtx.entranceMethodFullName) {
			// tracing finished, auto delete tracing
			// todo send events to facade ctx if necessary
			g.tracingGrIDMap.Delete(ctx.GrID)
		}
	}
}

func (g *GoRoutineTraceInterceptor) AddCurrentGRTracingContext(ctx *GoRoutineTracingContext) {
	g.tracingGrIDMap.Store(ctx.grID, ctx)
}

func (g *GoRoutineTraceInterceptor) DeleteCurrentGRTracingContext() {
	grID := goid.Get()
	g.tracingGrIDMap.Delete(grID)
}

// GetCurrentGRTracingContext return expected ctx type tracing ctx, FIXME: we now only support one facade tracing ctx type now
// e.g. if log debug gr is watching, another gr watching like 'trace' watching is invalid
func (g *GoRoutineTraceInterceptor) GetCurrentGRTracingContext(ctxType string) *GoRoutineTracingContext {
	grID := goid.Get()
	val, ok := g.tracingGrIDMap.Load(grID)
	if !ok {
		return nil
	}
	grTracingCtx := val.(*GoRoutineTracingContext)

	if grTracingCtx.GetFacadeCtx().Type() != ctxType {
		return nil
	}
	return grTracingCtx
}
