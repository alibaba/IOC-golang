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

package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		client, err := debugServiceClient.Watch(context.Background(), &boot.WatchRequest{
			InterfaceName:      args[0],
			ImplementationName: args[1],
			Method:             args[2],
			Input:              true,
			Output:             true,
		})
		if err != nil {
			panic(err)
		}
		for {
			msg, _ := client.Recv()
			fmt.Println()
			onToPrint := "Call"
			paramOrResponse := "Param"
			if !msg.IsParam {
				onToPrint = "Response"
				paramOrResponse = "Response"
			}
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.(%s).%s()", msg.InterfaceName, msg.ImplementationName, msg.MethodName)
			for index, p := range msg.GetParams() {
				color.Cyan("%s %d: %s", paramOrResponse, index+1, p)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watch)
}
