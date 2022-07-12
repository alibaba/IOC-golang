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
	log2 "github.com/opentracing/opentracing-go/log"
)

type traceLoggerWriter struct {
	spanGetter CurrentSpanGetter
	logKey     string
}

func (w *traceLoggerWriter) Write(p []byte) {
	currentSpan := w.spanGetter()
	if currentSpan != nil {
		// only when tracing, get span and write log
		currentSpan.LogFields(log2.String(w.logKey, string(p)))
	}
}

func newTraceLoggerWriter(spanGetter CurrentSpanGetter, logKey string) Writer {
	return &traceLoggerWriter{
		spanGetter: spanGetter,
		logKey:     logKey,
	}
}
