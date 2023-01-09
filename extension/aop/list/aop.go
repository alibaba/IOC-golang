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

package list

import (
	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/aop/common"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/extension/aop/list/api/ioc_golang/aop/list"
)

func init() {
	aop.RegisterAOP(aop.AOP{
		Name: "list",
		GRPCServiceRegister: func(server *grpc.Server) {
			listServiceImplSingleton, _ := GetlistServiceImplSingleton(nil)
			list.RegisterListServiceServer(server, listServiceImplSingleton)
		},
		ConfigLoader: func(aopConfig *common.Config) {
			_, _ = GetlistServiceImplSingleton(&listServiceImplParam{
				AllInterfaceMetadataMap: aop.GetAllInterfaceMetadata(),
				AppName:                 aopConfig.AppName,
			})
		},
	})
}
