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

package cli

import (
	"context"
	"fmt"
	"math"

	"github.com/alibaba/ioc-golang/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	aopLog "github.com/alibaba/ioc-golang/extension/aop/log"
	watchPB "github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/spf13/cobra"
)

func getWatchServiceClent(addr string) watchPB.WatchServiceClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	if err != nil {
		panic(err)
	}
	return watchPB.NewWatchServiceClient(conn)
}

var (
	debugHost string
	debugPort int
	maxDepth  int
	maxLength int
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logger.Red("invalid arguments, usage: iocli watch ${StructID} ${method}")
			return
		}
		debugServerAddr := fmt.Sprintf("%s:%d", debugHost, debugPort)
		debugServiceClient := getWatchServiceClent(debugServerAddr)
		logger.Cyan("iocli watch started, try to connect to debug server at %s", debugServerAddr)
		client, err := debugServiceClient.Watch(context.Background(), &watchPB.WatchRequest{
			Sdid:      args[0],
			MaxDepth:  int64(maxDepth),
			MaxLength: int64(maxLength),
			Method:    args[1],
		})
		if err != nil {
			panic(err)
		}
		logger.Cyan("debug server connected, watch info would be printed when invocation occurs")
		invocationCtxLogsGenerator, _ := aopLog.GetInvocationCtxLogsGeneratorSingleton()
		for {
			msg, err := client.Recv()
			if err != nil {
				logger.Red(err.Error())
				return
			}
			logger.Red(invocationCtxLogsGenerator.GetFunctionSignatureLogs(msg.Sdid, msg.MethodName, true))
			logger.Blue(invocationCtxLogsGenerator.GetParamsLogs(msg.GetParams(), true) + "\n")
			logger.Red(invocationCtxLogsGenerator.GetFunctionSignatureLogs(msg.Sdid, msg.MethodName, false))
			logger.Blue(invocationCtxLogsGenerator.GetParamsLogs(msg.GetReturnValues(), false) + "\n")
		}
	},
}

func init() {
	root.Cmd.AddCommand(watch)
	watch.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	watch.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	watch.Flags().IntVarP(&maxLength, "maxLength", "", 1000, "param value detail max length")
	watch.Flags().IntVarP(&maxDepth, "maxDepth", "", 5, "param value detail max depth")
}
