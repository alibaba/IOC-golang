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
	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/extension/autowire/allimpls"

	"github.com/alibaba/ioc-golang/logger"
)

type CurrentSpanGetter func() opentracing.Span

type Writer interface {
	Write(p []byte)
}

type TraceExtensionWriter interface {
	SetTraceLoggerWriter(writer Writer)
	Name() string
}

func RunRegisteredTraceLoggerWriterFunc(getter CurrentSpanGetter) {
	allTraceExtensionWriterImpls, err := allimpls.GetImpl(util.GetSDIDByStructPtr(new(TraceExtensionWriter)))
	if err != nil {
		logger.Red("[AOP Trace] Get all trace logger writer failed with error = %s", err)
		return
	}
	for _, impl := range allTraceExtensionWriterImpls.([]TraceExtensionWriter) {
		impl.SetTraceLoggerWriter(newTraceLoggerWriter(getter, impl.Name()))
		logger.Blue("[AOP] [Trace] Set trace logger to %s success", impl.Name())
	}
}
