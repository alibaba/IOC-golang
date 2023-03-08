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
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alibaba/ioc-golang/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alibaba/ioc-golang/extension/aop/trace/common"

	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/spf13/cobra"
)

func getTraceServiceClient(addr string) tracePB.TraceServiceClient {
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
		debugServiceClient := getTraceServiceClient(debugServerAddr)
		logger.Cyan("iocli trace started, try to connect to debug server at %s", debugServerAddr)
		client, err := debugServiceClient.Trace(context.Background(), &tracePB.TraceRequest{
			Sdid:                   sdid,
			Method:                 method,
			PushToCollectorAddress: pushToAddr,
			MaxDepth:               int64(maxDepth),
			MaxLength:              int64(maxLength),
		})
		if err != nil {
			panic(err)
		}
		logger.Cyan("debug server connected, tracing info would be printed every 5s (default)")

		jaegerCollectorEndpoint := common.GetJaegerCollectorEndpoint(pushToAddr)

		if pushToAddr != "" {
			logger.Cyan("try to push span batch data to %s", pushToAddr)
		}

		cacheData := bytes.Buffer{}
		if storeToFile != "" {
			logger.Cyan("Spans data is collecting, in order to save to %s", storeToFile)
			go func() {
				signals := make(chan os.Signal, 1)
				signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
				<-signals
				writeSpans(cacheData)
			}()
		}

		for {
			msg, err := client.Recv()
			if err != nil {
				logger.Red(err.Error())
				writeSpans(cacheData)
				return
			}
			for _, t := range msg.Traces {
				logger.Red("==================== Trace ====================")
				for _, span := range t.Spans {
					logger.Blue("Duration %dus, OperationName: %s, StartTime: %s, ReferenceSpans: %+v", span.GetDuration().Microseconds(), span.GetOperationName(), span.GetStartTime().Format("2006/01/02 15:04:05"), span.GetReferences())
					logger.Blue("====================")
				}
			}
			if data := msg.ThriftSerializedSpans; pushToAddr != "" && data != nil && len(data) > 0 {
				// try to push spans to collector
				body := bytes.NewBuffer(data)
				req, err := http.NewRequest("POST", jaegerCollectorEndpoint, body)
				if err != nil {
					logger.Red("New http request with url %s failed with error %s, ", jaegerCollectorEndpoint, err)
					continue
				}
				req.Header.Set("Content-Type", "application/x-thrift")
				go func() {
					// async post to collector
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						logger.Red("Http request with url %s failed with error %s, ", jaegerCollectorEndpoint, err)
						return
					}
					if resp.StatusCode >= http.StatusBadRequest {
						logger.Red(fmt.Sprintf("error from collector: %d", resp.StatusCode))
						return
					}
				}()
				cacheData.Write(data)
			}
		}
	},
}

func writeSpans(cacheData bytes.Buffer) {
	if err := os.WriteFile(storeToFile, cacheData.Bytes(), fs.ModePerm); err != nil {
		logger.Red("Write cached spans data to %s failed, error is %s", storeToFile, err.Error())
		os.Exit(1)
	}
	logger.Cyan("Write spans to %s finished", storeToFile)
}

var (
	debugHost   string
	debugPort   int
	maxDepth    int
	maxLength   int
	pushToAddr  string
	storeToFile string
)

func init() {
	root.Cmd.AddCommand(trace)
	trace.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	trace.Flags().IntVarP(&maxDepth, "maxDepth", "", 5, "param value detail max depth")
	trace.Flags().IntVarP(&maxLength, "maxLength", "", 1000, "param value detail max length")
	trace.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	trace.Flags().StringVar(&pushToAddr, "pushAddr", "", "push to jaeger collector address")
	trace.Flags().StringVar(&storeToFile, "store", "", "spans data store to file name")
}
