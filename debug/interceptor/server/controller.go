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

package server

import (
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/common"
	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/debug"

	debugCommon "github.com/alibaba/ioc-golang/debug/common"
	"github.com/alibaba/ioc-golang/debug/interceptor/trace"
	"github.com/alibaba/ioc-golang/debug/interceptor/watch"
)

func Start(debugConfig *debugCommon.Config, allInterfaceMetadataMap map[string]*debugCommon.StructMetadata) error {
	if debugConfig.AppName != "" {
		trace.SetAppName(debugConfig.AppName)
	}
	if collectorAddr := debugConfig.InterceptorsConfig.Trace.CollectorAddress; collectorAddr != "" {
		trace.SetCollectorAddress(collectorAddr)
	}
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&debug.DebugService_ServiceDesc, &DebugServerImpl{
		watchInterceptor:        watch.GetWatchInterceptor(),
		traceInterceptor:        trace.GetTraceInterceptor(),
		allInterfaceMetadataMap: allInterfaceMetadataMap,
	})

	lst, err := common.GetTCPListener(debugConfig.Port)
	if err != nil {
		return err
	}

	go func() {
		color.Blue("[Debug] Debug server listening at :%d", lst.Addr().(*net.TCPAddr).Port)
		if err := grpcServer.Serve(lst); err != nil {
			color.Red("[Debug] Debug server run with error = ", err)
			return
		}
	}()
	return nil
}
