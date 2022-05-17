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

package interceptor

import (
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	"github.com/alibaba/IOC-Golang/debug/api/ioc_golang/boot"
	"github.com/alibaba/IOC-Golang/debug/common"
)

func Start(port string, allInterfaceMetadataMap map[string]*common.DebugMetadata) error {
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&boot.DebugService_ServiceDesc, &DebugServerImpl{
		editInterceptor:         GetEditInterceptor(),
		watchInterceptor:        GetWatchInterceptor(),
		allInterfaceMetadataMap: allInterfaceMetadataMap,
	})
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		color.Red("[Debug] Debug server listening port :%s failed with error = %s", port, err)
		return err
	}
	go func() {
		color.Blue("[Debug] Debug server listening at :%s", port)
		if err := grpcServer.Serve(lst); err != nil {
			color.Red("[Debug] Debug server run with error = ", err)
			return
		}
	}()
	return nil
}
