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

package extension

import (
	"io"
	"log"

	traceLog "github.com/alibaba/ioc-golang/extension/aop/trace/log"
)

const logName = "log"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:constructFunc=newLogWriter
// +ioc:autowire:proxy:autoInjection=false
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/extension/aop/trace/log.TraceExtensionWriter

type logWriter struct {
	rawWriter         io.Writer
	traceLoggerWriter traceLog.Writer
}

func newLogWriter(l *logWriter) (*logWriter, error) {
	l.rawWriter = log.Writer()
	log.SetOutput(l)
	return l, nil
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	if l.traceLoggerWriter != nil {
		l.traceLoggerWriter.Write(p)
	}
	return l.rawWriter.Write(p)
}

func (l *logWriter) SetTraceLoggerWriter(traceLoggerWriter traceLog.Writer) {
	l.traceLoggerWriter = traceLoggerWriter
}

func (l *logWriter) Name() string {
	return logName
}

var _ traceLog.TraceExtensionWriter = &logWriter{}
