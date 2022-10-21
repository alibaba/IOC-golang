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
	"github.com/sirupsen/logrus"

	"github.com/alibaba/ioc-golang/aop/common"
	logPB "github.com/alibaba/ioc-golang/extension/aop/log/api/ioc_golang/aop/log"
	"github.com/alibaba/ioc-golang/logger"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

type logServiceImpl struct {
	logPB.UnimplementedLogServiceServer

	LogInterceptor logInterceptorIOCInterface `singleton:""`
}

func (l *logServiceImpl) Log(req *logPB.LogRequest, logServer logPB.LogService_LogServer) error {
	logger.Red("[Debug Server] Receive log request %+v\n", req.String())
	defer logger.Red("[Debug Server] Log request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethodName()
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMachers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	originLevel := logrus.GetLevel()
	originInvocationCtxLevel := l.LogInterceptor.GetInvocationCtxLogger().GetLevel()
	defer func() {
		logrus.SetLevel(originLevel)
		l.LogInterceptor.GetInvocationCtxLogger().SetLevel(originInvocationCtxLevel)
	}()

	targetLevel := logrus.DebugLevel
	if req.Level != 0 {
		targetLevel = logrus.Level(req.Level)
	}
	if req.Invocation {
		// change log level to debug if invocation ctx is enable
		logrus.SetLevel(logrus.DebugLevel)
		l.LogInterceptor.GetInvocationCtxLogger().SetLevel(logrus.DebugLevel)
	} else {
		// change log level to target level
		logrus.SetLevel(targetLevel)
		l.LogInterceptor.GetInvocationCtxLogger().SetLevel(targetLevel)
	}

	responseCh := make(chan *logPB.LogResponse, 100)

	logDebugCtx, _ := GetdebugLogContext(&debugLogContextParam{
		sdid:                sdid,
		methodName:          method,
		fieldMatcher:        fieldMatcher,
		ch:                  responseCh,
		invocationCtxEnable: req.Invocation,
		level:               targetLevel,
	})
	l.LogInterceptor.WatchLogs(logDebugCtx)

	done := logServer.Context().Done()
	go func() {
		for {
			select {
			case <-done:
				return
			case logResponse := <-responseCh:
				if logResponse != nil {
					_ = logServer.Send(logResponse)
				}
			}
		}
	}()
	<-done
	l.LogInterceptor.StopWatch()
	return nil
}
