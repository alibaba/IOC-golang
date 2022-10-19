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
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
	logPB "github.com/alibaba/ioc-golang/extension/aop/log/api/ioc_golang/aop/log"
	"github.com/alibaba/ioc-golang/extension/aop/trace/goroutine_trace"
)

// logInterceptor has two features
// [Feature1] print invocation logs to logrus default logger using debug level (default)
// [Feature2] support debug server logs write back and tracing logs write back of current goroutine of target method span

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy:autoInjection=false
// +ioc:autowire:paramType=logInterceptorParams
// +ioc:autowire:constructFunc=initLogInterceptor

type logInterceptor struct {
	invocationAOPLogFFunction  func(format string, args ...interface{})
	disable                    bool
	disablePrintParams         bool
	printParamsMaxDepth        int
	printParamsMaxLength       int
	InvocationCtxLogsGenerator *InvocationCtxLogsGenerator `singleton:""`

	GoRoutineInterceptor goroutine_trace.GoRoutineTraceInterceptorIOCInterface `singleton:""`

	logDebugCtx *debugLogContext
}

type logInterceptorParams struct {
	InvocationAOPLogConfig
}

func (p *logInterceptorParams) initLogInterceptor(interceptor *logInterceptor) (*logInterceptor, error) {
	if p.Disable {
		interceptor.disable = true
		return interceptor, nil
	}
	// log interceptor is enabled
	interceptor.disable = false

	// init logger
	if level, err := logrus.ParseLevel(p.Level); err != nil {
		return interceptor, err
	} else {
		switch level {
		case logrus.DebugLevel:
			interceptor.invocationAOPLogFFunction = logrus.Debugf
		case logrus.InfoLevel:
			interceptor.invocationAOPLogFFunction = logrus.Infof
		case logrus.WarnLevel:
			interceptor.invocationAOPLogFFunction = logrus.Warnf
		case logrus.ErrorLevel:
			interceptor.invocationAOPLogFFunction = logrus.Errorf
		case logrus.FatalLevel:
			interceptor.invocationAOPLogFFunction = logrus.Fatalf
		case logrus.PanicLevel:
			interceptor.invocationAOPLogFFunction = logrus.Panicf
		case logrus.TraceLevel:
			interceptor.invocationAOPLogFFunction = logrus.Tracef
		default:
			return interceptor, fmt.Errorf("invalid log level %d", level)
		}
	}

	interceptor.disablePrintParams = p.DisablePrintParams
	interceptor.printParamsMaxDepth = p.PrintParamsMaxDepth
	interceptor.printParamsMaxLength = p.PrintParamsMaxLength
	return interceptor, nil
}

func (w *logInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {
	// [Feature1]
	w.invocationAOPLogFFunction("[AOP Function Call] %s\n", w.InvocationCtxLogsGenerator.GetFunctionSignatureLogs(ctx.SDID, ctx.MethodName, true))
	w.invocationAOPLogFFunction("[AOP Param] %s\n\n", w.InvocationCtxLogsGenerator.GetParamsLogs(common.ReflectValues2Strings(ctx.Params, w.printParamsMaxDepth, w.printParamsMaxLength), true))

	// [Feature2]
	// 1. find if already in goroutine tracing
	if w.GoRoutineInterceptor.GetCurrentGRTracingContext(logGoRoutineInterceptorFacadeCtxType) != nil {
		w.GoRoutineInterceptor.BeforeInvoke(ctx)
		return
	}
	// current invocation not in goroutine tracing

	// 2. try to get matched method tracing context
	debugServerLogCtx := w.logDebugCtx
	if debugServerLogCtx == nil {
		return
	}
	// method tracing found,
	if debugServerLogCtx.fieldMatcher != nil && !debugServerLogCtx.fieldMatcher.Match(ctx.Params) {
		//doesn't match trace by method
		return
	}
	if debugServerLogCtx.sdid != "" && debugServerLogCtx.sdid != ctx.SDID {
		// doesn't match sdid
		return
	}
	if debugServerLogCtx.methodName != "" && debugServerLogCtx.methodName != ctx.MethodName {
		// doesn't match method
		return
	}
	// match method tracing context found

	// 3.start goroutine tracing
	// create facade ctx
	facadeLogCtx, err := GetlogGoRoutineInterceptorFacadeCtx(&logGoRoutineInterceptorFacadeCtxParam{
		logCh: debugServerLogCtx.ch,
		//traceEnable:            debugServerLogCtx.traceEnable,
		//entranceMethodFullName: ctx.MethodFullName,
	})
	if err != nil {
		log.Printf("logInterceptor GetlogGoRoutineInterceptorFacadeCtx failed with erorr = %s\n", err.Error())
		return
	}
	// create gr trace ctx
	grCtx, _ := goroutine_trace.GetGoRoutineTracingContext(&goroutine_trace.GoRoutineTracingContextParams{
		FacadeCtx:              facadeLogCtx,
		EntranceMethodFullName: ctx.MethodFullName,
	})

	// start tracing
	w.GoRoutineInterceptor.AddCurrentGRTracingContext(grCtx)
	w.GoRoutineInterceptor.BeforeInvoke(ctx)
}

func (w *logInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	// [Feature2]
	w.GoRoutineInterceptor.AfterInvoke(ctx)

	// [Feature1]
	w.invocationAOPLogFFunction("[AOP Function Response] %s\n", w.InvocationCtxLogsGenerator.GetFunctionSignatureLogs(ctx.SDID, ctx.MethodName, false))
	w.invocationAOPLogFFunction("[AOP Return Values] %s\n\n", w.InvocationCtxLogsGenerator.GetParamsLogs(common.ReflectValues2Strings(ctx.ReturnValues, w.printParamsMaxDepth, w.printParamsMaxLength), false))

}

func (w *logInterceptor) WatchLogs(logCtx *debugLogContext) {
	w.logDebugCtx = logCtx
}

func (w *logInterceptor) StopWatch() {
	w.logDebugCtx = nil
}

func (w *logInterceptor) NotifyLogs(content string) {
	//  get current gr span
	if grCtx := w.GoRoutineInterceptor.GetCurrentGRTracingContext(logGoRoutineInterceptorFacadeCtxType); grCtx != nil {
		// found current log gr span
		grLogFacadeCtx := grCtx.GetFacadeCtx().(*logGoRoutineInterceptorFacadeCtx)
		// todo: trace log switch
		// if grLogFacadeCtx.traceEnable || common.IsTraceEntrance(grLogFacadeCtx.entranceMethodFullName) {
		grLogFacadeCtx.pushContent(content)
	}
}

const logGoRoutineInterceptorFacadeCtxType = "logGoRoutineInterceptorFacadeCtx"

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=logGoRoutineInterceptorFacadeCtxParam
// +ioc:autowire:constructFunc=initLogGoRoutineInterceptorFacadeCtx

type logGoRoutineInterceptorFacadeCtx struct {
	logGoRoutineInterceptorFacadeCtxParam
}

type logGoRoutineInterceptorFacadeCtxParam struct {
	logCh chan *logPB.LogResponse
	//traceEnable            bool
	entranceMethodFullName string
}

func (p *logGoRoutineInterceptorFacadeCtxParam) initLogGoRoutineInterceptorFacadeCtx(l *logGoRoutineInterceptorFacadeCtx) (*logGoRoutineInterceptorFacadeCtx, error) {
	l.logGoRoutineInterceptorFacadeCtxParam = *p
	return l, nil
}
func (l *logGoRoutineInterceptorFacadeCtx) pushContent(content string) {
	if ch := l.logCh; ch != nil {
		select {
		case l.logCh <- &logPB.LogResponse{
			Content: content,
		}:
		default:
			log.Printf("[Log AOP] failed to write back content to debug server, %s\n", content)
		}
	}
}

func (l *logGoRoutineInterceptorFacadeCtx) BeforeInvoke(ctx *aop.InvocationContext) {
}

func (l *logGoRoutineInterceptorFacadeCtx) AfterInvoke(ctx *aop.InvocationContext) {
}

func (t *logGoRoutineInterceptorFacadeCtx) Type() string {
	return logGoRoutineInterceptorFacadeCtxType
}
