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

	"github.com/alibaba/ioc-golang/aop/api/ioc_golang/debug"

	"github.com/petermattis/goid"

	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/aop/interceptor/common"
)

type methodTracingContext struct {
	methodName   string
	sdid         string
	ch           chan *debug.TraceResponse // not useless
	fieldMatcher *common.FieldMatcher
	tracesMap    sync.Map // goroutine-id -> *goRoutineTracingContext
}

func NewTraceByMethodContext(sdid, method string, ch chan *debug.TraceResponse, fieldMatcher *common.FieldMatcher) *methodTracingContext {
	return &methodTracingContext{
		sdid:         sdid,
		methodName:   method,
		ch:           ch,
		fieldMatcher: fieldMatcher,
	}
}

func (t *methodTracingContext) addGoroutineTraceContext(grCtx *goRoutineTracingContext) {
	grID := goid.Get()
	t.tracesMap.Store(grID, grCtx)
}

type goRoutineTracingContext struct {
	MethodName        string
	grID              int64
	ClientSpanContext opentracing.SpanContext
	trace             *trace
}

func newGoRoutineTracingContextWithClientSpan(entranceMethod string, clientSpan opentracing.SpanContext) *goRoutineTracingContext {
	grID := goid.Get()
	return &goRoutineTracingContext{
		trace:             newTraceWithClientSpanContext(grID, entranceMethod, clientSpan),
		grID:              grID,
		ClientSpanContext: clientSpan,
	}
}

func newGoRoutineTracingContext(entranceMethod string) *goRoutineTracingContext {
	grID := goid.Get()
	return &goRoutineTracingContext{
		trace: newTrace(grID, entranceMethod),
		grID:  grID,
	}
}

func (g *goRoutineTracingContext) getTrace() *trace {
	return g.trace
}
