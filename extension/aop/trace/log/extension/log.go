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

func init() {
	traceLog.RegisterTraceLoggerWriterFunc(logName, func(traceLoggerWrtter traceLog.Writer) {
		log.SetOutput(getLogWriter(log.Writer(), traceLoggerWrtter))
	})
}

type logWriter struct {
	rawWriter         io.Writer
	traceLoggerWriter traceLog.Writer
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.traceLoggerWriter.Write(p)
	return l.rawWriter.Write(p)
}

func getLogWriter(rawWriter io.Writer, traceLoggerWriter traceLog.Writer) io.Writer {
	return &logWriter{
		rawWriter:         rawWriter,
		traceLoggerWriter: traceLoggerWriter,
	}
}
