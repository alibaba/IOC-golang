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
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type trace struct {
	grID           int64
	currentSpan    *spanWithParent
	entranceMethod string
}

func newTraceWithClientSpanContext(grID int64, entranceMethod string, clientSpanContext opentracing.SpanContext) *trace {
	rootSpan := getGlobalTracer().getRawTracer().StartSpan(entranceMethod, ext.RPCServerOption(clientSpanContext))
	return &trace{
		grID:           grID,
		currentSpan:    newSpanWithParent(rootSpan, nil),
		entranceMethod: entranceMethod,
	}
}

func newTrace(grID int64, entranceMethod string) *trace {
	rootSpan := getGlobalTracer().getRawTracer().StartSpan(entranceMethod)
	return &trace{
		grID:           grID,
		currentSpan:    newSpanWithParent(rootSpan, nil),
		entranceMethod: entranceMethod,
	}
}

func (t *trace) addChildSpan(name string) *spanWithParent {
	func1Span := getGlobalTracer().getRawTracer().StartSpan(name, opentracing.ChildOf(t.currentSpan.span.Context()), opentracing.StartTime(time.Now()))
	innerChildSpan := newSpanWithParent(func1Span, t.currentSpan)
	t.currentSpan = innerChildSpan
	return t.currentSpan
}

func (t *trace) returnSpan() {
	t.currentSpan.span.Finish()
	t.currentSpan = t.currentSpan.parentSpan
}
