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

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/debug"

	"github.com/spf13/cobra"
)

var watch = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(fmt.Sprintf("%s:%d", debugHost, debugPort))
		client, err := debugServiceClient.Watch(context.Background(), &debug.WatchRequest{
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
			paramOrResponse := "Param"
			onToPrint := "Call"
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.%s()", msg.Sdid, msg.MethodName)
			for index, p := range msg.GetParams() {
				color.Cyan("%s %d: %s", paramOrResponse, index+1, p)
			}

			onToPrint = "Response"
			paramOrResponse = "Response"
			color.Red("========== On %s ==========\n", onToPrint)
			color.Red("%s.%s()", msg.Sdid, msg.MethodName)
			for index, p := range msg.GetReturnValues() {
				color.Cyan("%s %d: %s", paramOrResponse, index+1, p)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watch)
	watch.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	watch.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
