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

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	logPB "github.com/alibaba/ioc-golang/extension/aop/log/api/ioc_golang/aop/log"
	"github.com/alibaba/ioc-golang/iocli/root"
	"github.com/alibaba/ioc-golang/logger"
)

func getLogServiceClient(addr string) logPB.LogServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return logPB.NewLogServiceClient(conn)
}

var call = &cobra.Command{
	Use:     "log",
	Aliases: []string{"logs"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			logger.Red("iocli log command needs 3 arguments: \n${autowireType} ${sdid} ${methodName} \n")
			return
		}
		autowireType := args[0]
		sdid := args[1]
		methodName := args[2]
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			logger.Red("level %s is invalid", level)
			return
		}

		if printInvocationCtx {
			logger.Cyan("print invocation ctx")
		}

		callServiceClient := getLogServiceClient(fmt.Sprintf("%s:%d", debugHost, debugPort))
		logSvcClient, err := callServiceClient.Log(context.Background(), &logPB.LogRequest{
			Sdid:         sdid,
			MethodName:   methodName,
			AutowireType: autowireType,
			Level:        int64(logrusLevel),
			Invocation:   printInvocationCtx,
		})
		if err != nil {
			logger.Red(err.Error())
			return
		}

		for {
			rsp, err := logSvcClient.Recv()
			if err != nil {
				logger.Red(err.Error())
				return
			}
			logger.Blue(rsp.Content)
		}
	},
}

var (
	debugHost          string
	debugPort          int
	level              string
	printInvocationCtx bool
)

func init() {
	root.Cmd.AddCommand(call)
	call.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	call.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	call.Flags().StringVar(&level, "level", "debug", "pull logs level")
	call.Flags().BoolVar(&printInvocationCtx, "invocation", false, "if print invocation context info")
}
