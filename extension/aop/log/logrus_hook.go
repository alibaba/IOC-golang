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

package call

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/alibaba/ioc-golang/aop"
)

// LogrusIOCCtxHook
// [Feature1] append some keys as field of logs
// [Feature2] send log content to current gr ctx, to let log interceptor collect if necessary

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:loadAtOnce=true
// +ioc:autowire:constructFunc=newLogrusIOCCtxHook
// +ioc:autowire:proxy:autoInjection=false

type LogrusIOCCtxHook struct {
	originWriter io.Writer

	structIDKey   string
	methodNameKey string
	grIDKey       string

	LogInterceptor logInterceptorIOCInterface `singleton:""`
}

func newLogrusIOCCtxHook(l *LogrusIOCCtxHook) (*LogrusIOCCtxHook, error) {
	logrus.AddHook(l)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})
	l.structIDKey = "structID"
	l.methodNameKey = "methodName"
	l.grIDKey = "grID"
	return l, nil
}

func (l *LogrusIOCCtxHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (l *LogrusIOCCtxHook) Fire(entry *logrus.Entry) error {
	// [Feature1]
	aopInvocationCtx := aop.GetCurrentInvocationCtx()
	if aopInvocationCtx == nil {
		return nil
	}
	if entry.Data == nil {
		entry.Data = make(map[string]interface{})
	}
	entry.Data[l.structIDKey] = aopInvocationCtx.SDID
	entry.Data[l.methodNameKey] = aopInvocationCtx.MethodName
	entry.Data[l.grIDKey] = aopInvocationCtx.GrID

	// [Feature2]
	content, err := entry.String()
	if err != nil {
		return err
	}
	l.LogInterceptor.NotifyLogs(content)
	return nil
}

/*
	SetLogLevel is used to dynamic change log level using 'iocli call' command

 	iocli call singleton github.com/alibaba/ioc-golang/extension/aop/log.LogrusIOCCtxHook SetLogLevel --params "[2]"

params:
	PanicLevel: 0
	FatalLevel: 1
	ErrorLevel: 2
	WarnLevel: 3
	InfoLevel: 4
	DebugLevel: 5
	TraceLevel: 6
*/
func (l *LogrusIOCCtxHook) SetLogLevel(level uint32) {
	logrus.SetLevel(logrus.Level(level))
}
