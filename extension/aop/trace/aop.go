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

package trace

import (
	"fmt"

	"github.com/alibaba/ioc-golang/logger"

	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
	"github.com/alibaba/ioc-golang/extension/aop/trace/log"
	_ "github.com/alibaba/ioc-golang/extension/aop/trace/log/extension"
)

const Name = "trace"

func init() {
	aop.RegisterAOP(aop.AOP{
		Name: Name,
		InterceptorFactory: func() aop.Interceptor {
			interceptor, _ := GettraceInterceptorIOCInterfaceSingleton()
			// inject logger interceptor
			log.RunRegisteredTraceLoggerWriterFunc(interceptor.GetCurrentSpan)
			return interceptor
		},
		RPCInterceptorFactory: func() aop.RPCInterceptor {
			traceRPCInterceptor, err := GetrpcInterceptorIOCInterfaceSingleton()
			if err != nil {
				logger.Red("trace AOP GetrpcInterceptorIOCInterfaceSingleton failed with error = %s", err.Error())
			}
			return traceRPCInterceptor
		},
		GRPCServiceRegister: func(server *grpc.Server) {
			rpcService, _ := GettraceServiceImplSingleton()
			tracePB.RegisterTraceServiceServer(server, rpcService)
		},
		ConfigLoader: func(aopConfig *common.Config) {
			if aopConfig.AppName != "" {
				setAppName(aopConfig.AppName)
			}
			traceConfig := &TraceConfig{}
			if err := config.LoadConfigByPrefix(fmt.Sprintf("%s.%s", common.IOCGolangAOPConfigPrefix, Name), traceConfig); err == nil {
				// found property
				setCollectorAddress(traceConfig.CollectorAddress)
			}
			if traceConfig.ValueDepth != 0 {
				valueDepth = traceConfig.ValueDepth
			}
		},
	})
}
