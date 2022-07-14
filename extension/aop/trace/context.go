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
	"github.com/petermattis/goid"

	traceCommon "github.com/alibaba/ioc-golang/extension/aop/trace/common"

	"github.com/alibaba/ioc-golang/aop/common"
)

type methodTracingContext struct {
	methodName   string
	sdid         string
	fieldMatcher *common.FieldMatcher
	maxDepth     int64
	maxLength    int64
}

func newTraceByMethodContext(sdid, method string, fieldMatcher *common.FieldMatcher, maxDepth, maxLength int64) *methodTracingContext {
	return &methodTracingContext{
		sdid:         sdid,
		methodName:   method,
		fieldMatcher: fieldMatcher,
		maxLength:    maxLength,
		maxDepth:     maxDepth,
	}
}

type goRoutineTracingContext struct {
	MethodName        string
	grID              int64
	ClientSpanContext opentracing.SpanContext
	trace             *trace
	maxDepth          int64
	maxLength         int64
}

func newGoRoutineTracingContextWithClientSpan(entranceMethod string, clientSpan opentracing.SpanContext) *goRoutineTracingContext {
	grID := goid.Get()
	return &goRoutineTracingContext{
		trace:             newTraceWithClientSpanContext(grID, entranceMethod, clientSpan),
		grID:              grID,
		ClientSpanContext: clientSpan,
		maxDepth:          traceCommon.DefaultRecordValuesDepth,
		maxLength:         traceCommon.DefaultRecordValuesLength,
	}
}

func newGoRoutineTracingContext(entranceMethod string, maxDepth, maxLength int64) *goRoutineTracingContext {
	grID := goid.Get()
	return &goRoutineTracingContext{
		trace:     newTrace(grID, entranceMethod),
		grID:      grID,
		maxDepth:  maxDepth,
		maxLength: maxLength,
	}
}

func (g *goRoutineTracingContext) getTrace() *trace {
	return g.trace
}
