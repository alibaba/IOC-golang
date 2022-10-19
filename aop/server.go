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

package aop

import (
	"math"
	"net"

	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/logger"
)

func start(debugConfig *common.Config) error {
	for _, cl := range configLoaderFuncs {
		cl(debugConfig)
	}

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(math.MaxInt32))
	for _, register := range grpcServiceRegisters {
		register(grpcServer)
	}

	lst, err := common.GetTCPListener(debugConfig.DebugServer.Port)
	if err != nil {
		return err
	}

	go func() {
		logger.Blue("[Debug] Debug server listening at :%d", lst.Addr().(*net.TCPAddr).Port)
		if err := grpcServer.Serve(lst); err != nil {
			logger.Red("[Debug] Debug server run with error = ", err)
			return
		}
	}()
	return nil
}
