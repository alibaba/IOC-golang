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

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/debug"
)

var trace = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(fmt.Sprintf("%s:%d", debugHost, debugPort))
		client, err := debugServiceClient.Trace(context.Background(), &debug.TraceRequest{
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

func init() {
	rootCmd.AddCommand(trace)
	trace.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	trace.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
}
