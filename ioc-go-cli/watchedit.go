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
	"log"

	"github.com/spf13/cobra"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

// watchEdit todo
var watchEdit = &cobra.Command{
	Use: "watchEdit",
	Run: func(cmd *cobra.Command, args []string) {
		debugServiceClient := getDebugServiceClent(defaultDebugAddr)
		watchEditClient, err := debugServiceClient.WatchEdit(context.Background())
		if err != nil {
			panic(err)
		}
		if err := watchEditClient.Send(&boot.WatchEditRequest{
			InterfaceName:      args[0],
			ImplementationName: args[1],
			Method:             args[2],
			IsEdit:             false,
		}); err != nil {
			panic(err)
		}
		for {
			rsp, err := watchEditClient.Recv()
			if err != nil {
				log.Printf("recv error = %s\n", err)
				return
			}
			fmt.Println(rsp.Params)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchEdit)
}
