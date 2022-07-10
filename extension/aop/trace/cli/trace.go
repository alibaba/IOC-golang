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
	"bytes"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alibaba/ioc-golang/extension/aop/trace/common"

	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func getTraceServiceClent(addr string) tracePB.TraceServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return tracePB.NewTraceServiceClient(conn)
}

var trace = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		sdid := ""
		method := ""
		if len(args) > 0 {
			sdid = args[0]
		}
		if len(args) > 1 {
			method = args[1]
		}
		debugServerAddr := fmt.Sprintf("%s:%d", debugHost, debugPort)
		debugServiceClient := getTraceServiceClent(debugServerAddr)
		color.Cyan("iocli trace started, try to connect to debug server at %s", debugServerAddr)
		client, err := debugServiceClient.Trace(context.Background(), &tracePB.TraceRequest{
			Sdid:                   sdid,
			Method:                 method,
			PushToCollectorAddress: pushToAddr,
		})
		if err != nil {
			panic(err)
		}
		color.Cyan("debug server connected, tracing info would be printed every 5s (default)")

		jaegerCollectorEndpoint := common.GetJaegerCollectorEndpoint(pushToAddr)

		if pushToAddr != "" {
			color.Cyan("try to push span batch data to %s", pushToAddr)
		}

		for {
			msg, err := client.Recv()
			if err != nil {
				color.Red(err.Error())
				return
			}
			for _, t := range msg.Traces {
				color.Red("==================== Trace ====================")
				for _, span := range t.Spans {
					color.Blue("Duration %dus, OperationName: %s, StartTime: %s, ReferenceSpans: %+v", span.GetDuration().Microseconds(), span.GetOperationName(), span.GetStartTime().Format("2006/01/02 15:04:05"), span.GetReferences())
					color.Blue("====================")
				}
			}
			if data := msg.ThriftSerializedSpans; pushToAddr != "" && data != nil && len(data) > 0 {
				body := bytes.NewBuffer(data)
				req, err := http.NewRequest("POST", jaegerCollectorEndpoint, body)
				if err != nil {
					color.Red("New http request with url %s failed with error %s, ", jaegerCollectorEndpoint, err)
					continue
				}
				req.Header.Set("Content-Type", "application/x-thrift")
				go func() {
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						color.Red("Http request with url %s failed with error %s, ", jaegerCollectorEndpoint, err)
						return
					}
					if resp.StatusCode >= http.StatusBadRequest {
						color.Red(fmt.Sprintf("error from collector: %d", resp.StatusCode))
						return
					}
				}()
			}
		}
	},
}

var (
	debugHost  string
	debugPort  int
	pushToAddr string
)

func init() {
	root.Cmd.AddCommand(trace)
	trace.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	trace.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	trace.Flags().StringVar(&pushToAddr, "pushAddr", "", "push to jaeger collector address")
}
