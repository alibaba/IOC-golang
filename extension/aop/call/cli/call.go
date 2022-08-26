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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	callPB "github.com/alibaba/ioc-golang/extension/aop/call/api/ioc_golang/aop/call"
	"github.com/alibaba/ioc-golang/iocli/root"
	"github.com/alibaba/ioc-golang/logger"
)

func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

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
		if len(args) < 3 {
			logger.Red("iocli call command needs 3 arguments: \n${autowireType} ${sdid} ${methodName} \n")
			return
		}
		autowireType := args[0]
		sdid := args[1]
		methodName := args[2]

		paramsJSON := params
		if paramsJSON == "" && paramsFile != "" {
			data, err := ioutil.ReadFile(paramsFile)
			if err != nil {
				logger.Red("iocli call command read param json from %s failed with error = %s", paramsFile, err.Error())
				return
			}
			paramsJSON = string(data)
		}

		if len(paramsJSON) > 0 && !isValidJSON(paramsJSON) {
			logger.Red("iocli call command param %s invalid: \n", paramsJSON)
			return
		}

		callServiceClient := getCallServiceClient(fmt.Sprintf("%s:%d", debugHost, debugPort))
		// todo
		rsp, err := callServiceClient.Call(context.Background(), &callPB.CallRequest{
			Sdid:         sdid,
			MethodName:   methodName,
			AutowireType: autowireType,
			Params:       []byte(paramsJSON),
		})
		if err != nil {
			logger.Red(err.Error())
			return
		}

		logger.Blue("Call %s: %s.%s() success!\nParam = %s\nReturn values = %s", autowireType, sdid, methodName, paramsJSON, string(rsp.ReturnValues))
	},
}

var (
	debugHost  string
	debugPort  int
	params     string
	paramsFile string
)

func init() {
	root.Cmd.AddCommand(call)
	call.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	call.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	call.Flags().StringVar(&params, "params", "", "request call params json string")
	call.Flags().StringVar(&paramsFile, "paramsFile", "", "request call params json file path")
}
