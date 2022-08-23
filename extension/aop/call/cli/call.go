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
	callPB "github.com/alibaba/ioc-golang/extension/aop/call/api/ioc_golang/aop/call"
	"github.com/alibaba/ioc-golang/iocli/root"
	"github.com/alibaba/ioc-golang/logger"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getCallServiceClient(addr string) callPB.CallServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return callPB.NewCallServiceClient(conn)
}

var call = &cobra.Command{
	Use: "call",
	Run: func(cmd *cobra.Command, args []string) {
		callServiceClient := getCallServiceClient(fmt.Sprintf("%s:%d", debugHost, debugPort))
		// todo
		_, err := callServiceClient.Call(context.Background(), nil)
		if err != nil {
			logger.Red(err.Error())
			return
		}
	},
}

var (
	debugHost string
	debugPort int
)

func init() {
	root.Cmd.AddCommand(call)
	call.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	call.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
