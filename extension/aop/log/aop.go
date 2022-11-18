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

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"

	logPB "github.com/alibaba/ioc-golang/extension/aop/log/api/ioc_golang/aop/log"

	"github.com/alibaba/ioc-golang/aop"
)

const Name = "log"

func init() {
	aop.RegisterAOP(aop.AOP{
		Name: Name,
		GRPCServiceRegister: func(server *grpc.Server) {
			logServiceImplSingleton, _ := GetlogServiceImplSingleton()
			logPB.RegisterLogServiceServer(server, logServiceImplSingleton)
		},
		InterceptorFactory: func() aop.Interceptor {
			// get loaded logInterceptor singleton
			logInterceptorSingleton, _ := GetlogInterceptorIOCInterfaceSingleton(nil)
			return logInterceptorSingleton
		},
		ConfigLoader: func(aopConfig *common.Config) {
			logConfig := &LogConfig{}
			_ = config.LoadConfigByPrefix(fmt.Sprintf("%s.%s", common.IOCGolangAOPConfigPrefix, Name), logConfig)
			logConfig.fillDefaultConfig()

			// init logInterceptor singleton
			_, _ = GetlogInterceptorIOCInterfaceSingleton(&logInterceptorParams{
				InvocationAOPLogConfig: logConfig.InvocationAOPLogConfig,
			})

			// init global logrus hook
			globalLogLevel, _ := logrus.ParseLevel(logConfig.Level)
			_, _ = GetGlobalLogrusIOCCtxHookIOCInterfaceSingleton(&globalLogrusIOCCtxHookParam{
				globalLogLevel:       globalLogLevel,
				globalLoggerReadOnly: logConfig.GlobalLoggerReadOnly,
			})
		},
	})
}
