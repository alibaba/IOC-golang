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
	"fmt"

	"github.com/inconshreveable/log15"

	traceLog "github.com/alibaba/ioc-golang/extension/aop/trace/log"
)

const log15Name = "log15"

func init() {
	traceLog.RegisterTraceLoggerWriterFunc(log15Name, func(traceLoggerWriter traceLog.Writer) {
		log15.Root().SetHandler(newLog15Handler(log15.Root().GetHandler(), traceLoggerWriter))
	})
}

type log15Handler struct {
	traceLoggerWriter traceLog.Writer
	handler           log15.Handler
}

func (l *log15Handler) Log(r *log15.Record) error {
	l.traceLoggerWriter.Write([]byte(fmt.Sprintf("%+v", *r)))
	return l.handler.Log(r)
}

func newLog15Handler(handler log15.Handler, traceLoggerWriter traceLog.Writer) log15.Handler {
	return &log15Handler{
		traceLoggerWriter: traceLoggerWriter,
		handler:           handler,
	}
}
