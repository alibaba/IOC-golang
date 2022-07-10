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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	watchPB "github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/fatih/color"

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
	depth     int
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			color.Red("invalid arguments, usage: iocli watch ${StructID} ${method}")
			return
		}
		debugServerAddr := fmt.Sprintf("%s:%d", debugHost, debugPort)
		debugServiceClient := getWatchServiceClent(debugServerAddr)
		color.Cyan("iocli watch started, try to connect to debug server at %s", debugServerAddr)
		client, err := debugServiceClient.Watch(context.Background(), &watchPB.WatchRequest{
			Sdid:     args[0],
			MaxDepth: int64(depth),
			Method:   args[1],
		})
		if err != nil {
			panic(err)
		}
		color.Cyan("debug server connected, watch info would be printed when invocation occurs, param info max depth = %d", depth)
		for {
			msg, err := client.Recv()
			if err != nil {
				color.Red(err.Error())
				return
			}
			paramOrResponse := "Param"
			onToPrint := "Call"
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.%s()", msg.Sdid, msg.MethodName)
			for index, p := range msg.GetParams() {
				color.Blue("%s %d: %s", paramOrResponse, index+1, p)
			}

			onToPrint = "Response"
			paramOrResponse = "Response"
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.%s()", msg.Sdid, msg.MethodName)
			for index, p := range msg.GetReturnValues() {
				color.Blue("%s %d: %s", paramOrResponse, index+1, p)
			}
		}
	},
}

func init() {
	root.Cmd.AddCommand(watch)
	watch.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	watch.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	watch.Flags().IntVarP(&depth, "depth", "d", 5, "value depth")
}
