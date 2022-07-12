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

package log

import (
	"github.com/fatih/color"
	"github.com/opentracing/opentracing-go"
)

type CurrentSpanGetter func() opentracing.Span

type SetTraceLoggerWriterFunc func(traceLoggerWriter Writer)

var rawLoggerWriterMap = make(map[string]SetTraceLoggerWriterFunc)

func RegisterTraceLoggerWriterFunc(name string, rawLoggerWriterManager SetTraceLoggerWriterFunc) {
	rawLoggerWriterMap[name] = rawLoggerWriterManager
}

func RunRegisteredTraceLoggerWriterFunc(getter CurrentSpanGetter) {
	for name, f := range rawLoggerWriterMap {
		f(newTraceLoggerWriter(getter, name))
		color.Blue("[AOP] [Trace] Set trace logger to %s success", name)
	}
}
