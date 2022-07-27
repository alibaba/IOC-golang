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

package watch

import (
	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/aop"
	watchPB "github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
)

const Name = "watch"

func init() {
	aop.RegisterAOP(aop.AOP{
		Name: Name,
		InterceptorFactory: func() aop.Interceptor {
			impl, _ := GetinterceptorImplIOCInterfaceSingleton()
			return impl
		},
		GRPCServiceRegister: func(server *grpc.Server) {
			watchInterface, _ := GetwatchServiceSingleton()
			watchPB.RegisterWatchServiceServer(server, watchInterface)
		},
	})
}
