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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alibaba/IOC-Golang/debug/api/ioc_golang/boot"
)

const (
	defaultDebugAddr = "localhost:1999"
)

var rootCmd = &cobra.Command{
	Use: "ioc-go-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func getDebugServiceClent(addr string) boot.DebugServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return boot.NewDebugServiceClient(conn)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}
