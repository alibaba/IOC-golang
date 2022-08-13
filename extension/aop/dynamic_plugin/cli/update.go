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
	"io/ioutil"
	"math"

	"github.com/alibaba/ioc-golang/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"

	dynamicPluginPB "github.com/alibaba/ioc-golang/extension/aop/dynamic_plugin/api/ioc_golang/aop/dynamic_plugin"
)

func getDynamicPluginServiceClient(addr string) dynamicPluginPB.DynamicPluginServiceClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	if err != nil {
		panic(err)
	}
	return dynamicPluginPB.NewDynamicPluginServiceClient(conn)
}

var (
	debugHost string
	debugPort int
)

var updateCommand = &cobra.Command{
	Use:     "update",
	Short:   "Update go plugin to application",
	Long:    "Update go plugin to application",
	Example: "iocli goplugin update singleton main.ServiceImpl1 ./service1.so Service1Plugin\nThe service1.go should be built by: 'go build -buildmode=plugin -o service1.so .'",
	Run: func(cmd *cobra.Command, args []string) {
		dynamicPluginServerAddr := fmt.Sprintf("%s:%d", debugHost, debugPort)
		dynamicPluginServiceClient := getDynamicPluginServiceClient(dynamicPluginServerAddr)
		if len(args) < 4 {
			logger.Red("iocli goplugin update command needs 4 arguments: \n${autowireType} ${sdid} ${pluginPath} ${pluginName} \nlike: \niocli goplugin update singleton main.ServiceImpl1 ./service1.so Service1Plugin\nThe service1.go should be built by: 'go build -buildmode=plugin -o service1.so .'")
			return
		}
		autowireType := args[0]
		sdid := args[1]
		pluginFilePath := args[2]
		pluginName := args[3]

		pluginFile, err := ioutil.ReadFile(pluginFilePath)
		if err != nil {
			logger.Red("Read plugin file %s failed with error = %s", pluginFilePath, err)
			return
		}

		logger.Cyan("Read plugin file %s success!\niocli goplugin update started, try to connect to debug server at %s\n"+
			"try to update goplugin with autowire type = %s, sdid = %s, pluginName = %s", pluginFilePath, dynamicPluginServerAddr, autowireType, sdid, pluginName)
		rsp, err := dynamicPluginServiceClient.Update(context.Background(), &dynamicPluginPB.DynamicPluginUpdateRequest{
			Sdid:         sdid,
			PluginFile:   pluginFile,
			PluginName:   pluginName,
			AutowireType: autowireType,
		})
		if err != nil {
			logger.Red("Update plugin failed with error = %s", err)
			return
		}

		if !rsp.Success {
			logger.Red("Update plugin failed with error = %s", rsp.GetSuccess())
			return
		}

		logger.Cyan("Update plugin success!")
	},
}

func init() {
	pluginCommand.AddCommand(updateCommand)

	updateCommand.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	updateCommand.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
