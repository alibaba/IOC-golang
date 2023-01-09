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

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	listPB "github.com/alibaba/ioc-golang/extension/aop/list/api/ioc_golang/aop/list"
	"github.com/alibaba/ioc-golang/iocli/root"
	"github.com/alibaba/ioc-golang/logger"
)

func getListServiceClent(addr string) listPB.ListServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return listPB.NewListServiceClient(conn)
}

var list = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		listServiceClient := getListServiceClent(fmt.Sprintf("%s:%d", debugHost, debugPort))
		rsp, err := listServiceClient.List(context.Background(), &emptypb.Empty{})
		if err != nil {
			logger.Red(err.Error())
			return
		}
		if rsp.AppName != "" {
			logger.Blue("appName: %s", rsp.GetAppName())
		}
		for _, v := range rsp.ServiceMetadata {
			logger.Blue(v.ImplementationName)
			logger.Blue("%s", v.Methods)
			logger.Blue("")
		}
	},
}

var (
	debugHost string
	debugPort int
)

func init() {
	root.Cmd.AddCommand(list)
	list.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	list.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
