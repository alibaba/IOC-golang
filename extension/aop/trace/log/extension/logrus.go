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

	"github.com/sirupsen/logrus"

	traceLog "github.com/alibaba/ioc-golang/extension/aop/trace/log"
)

const logrusName = "logrus"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:constructFunc=newLogrusWriter
// +ioc:autowire:proxy:autoInjection=false
// +ioc:autowire:implements=github.com/alibaba/ioc-golang/extension/aop/trace/log.TraceExtensionWriter

type logrusWriter struct {
	originWriter      io.Writer
	traceLoggerWriter traceLog.Writer
}

func newLogrusWriter(l *logrusWriter) (*logrusWriter, error) {
	l.originWriter = logrus.StandardLogger().Out
	logrus.SetOutput(l)
	return l, nil
}

func (l *logrusWriter) Write(p []byte) (n int, err error) {
	if l.traceLoggerWriter != nil {
		l.traceLoggerWriter.Write(p)
	}
	return l.originWriter.Write(p)
}

func (l *logrusWriter) SetTraceLoggerWriter(traceLoggerWriter traceLog.Writer) {
	l.traceLoggerWriter = traceLoggerWriter
}

func (l *logrusWriter) Name() string {
	return logrusName
}

var _ traceLog.TraceExtensionWriter = &logrusWriter{}
