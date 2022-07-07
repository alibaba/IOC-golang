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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func getTraceServiceClent(addr string) tracePB.TraceServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return tracePB.NewTraceServiceClient(conn)
}

var trace = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getTraceServiceClent(fmt.Sprintf("%s:%d", debugHost, debugPort))
		client, err := debugServiceClient.Trace(context.Background(), &tracePB.TraceRequest{
			Sdid:   args[0],
			Method: args[1],
		})
		if err != nil {
			panic(err)
		}
		for {
			msg, err := client.Recv()
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Blue("Tracing data is sending to %s", msg.CollectorAddress)
		}
	},
}

var (
	debugHost string
	debugPort int
)

func init() {
	root.Cmd.AddCommand(trace)
	trace.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	trace.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
