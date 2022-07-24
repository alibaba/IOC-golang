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
	"math"
	"os"
	"os/signal"
	"sort"
	"sync"
	"time"

	"github.com/alibaba/ioc-golang/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	monitorPB "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
	"github.com/alibaba/ioc-golang/iocli/root"

	"github.com/spf13/cobra"
)

func getMonitorServiceClient(addr string) monitorPB.MonitorServiceClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	if err != nil {
		panic(err)
	}
	return monitorPB.NewMonitorServiceClient(conn)
}

var (
	debugHost string
	debugPort int
	interval  int
)

var monitorCommand = &cobra.Command{
	Use: "monitor",
	Run: func(cmd *cobra.Command, args []string) {
		debugServerAddr := fmt.Sprintf("%s:%d", debugHost, debugPort)
		debugServiceClient := getMonitorServiceClient(debugServerAddr)
		sdid := ""
		method := ""
		if len(args) > 0 {
			sdid = args[0]
		}
		if len(args) > 1 {
			method = args[1]
		}
		logger.Cyan("iocli monitor started, try to connect to debug server at %s", debugServerAddr)
		client, err := debugServiceClient.Monitor(context.Background(), &monitorPB.MonitorRequest{
			Sdid:     sdid,
			Method:   method,
			Interval: int64(interval),
		})
		if err != nil {
			panic(err)
		}
		logger.Cyan("debug server connected, monitor info would be printed every %ds", interval)

		allMonitorResponseItemsMap := make(map[string][]*monitorPB.MonitorResponseItem, 0)
		allMonitorResponseItemsLock := sync.RWMutex{}
		startTime := time.Now().UnixMilli()
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt, os.Kill)
			select {
			case <-signals:
				fmt.Println()
				logger.Red("Got interrupt signal, collecting data during %dms", time.Now().UnixMilli()-startTime)
				logger.Red("====================Collection====================")
				logger.Red("%s", time.Now().Format("2006/01/02 15:04:05"))
				allMonitorResponseItemsLock.RLock()

				allMethods := make(methodSorter, 0)
				for k, _ := range allMonitorResponseItemsMap {
					allMethods = append(allMethods, k)
				}
				sort.Sort(allMethods)

				for _, methodKey := range allMethods {
					// get method all records
					items := allMonitorResponseItemsMap[methodKey]
					logger.Blue(methodKey)

					// init value
					total := int64(0)
					success := int64(0)
					fail := int64(0)
					avgRT := float32(0)
					avgFailRate := float32(0)

					allAvgRTS := make([]float32, 0)
					allFailRates := make([]float32, 0)

					// calculate total and average
					for _, item := range items {
						allAvgRTS = append(allAvgRTS, item.AvgRT)
						allFailRates = append(allFailRates, item.FailRate)

						total += item.Total
						fail += item.Fail
						success += item.Success
					}
					avgRT = getAverageFloat32(allAvgRTS)
					avgFailRate = getAverageFloat32(allFailRates)

					// print information
					logger.Blue(fmt.Sprintf("Total: %d, Success: %d, Fail: %d, AvgRT: %.2fus, FailRate: %.2f%%",
						total, success, fail, avgRT, avgFailRate))
				}

				allMonitorResponseItemsLock.RUnlock()
				// shutdown
				os.Exit(0)
			}
		}()
		for {
			msg, err := client.Recv()
			if err != nil {
				logger.Red(err.Error())
				return
			}
			logger.Red("====================")
			logger.Red("%s", time.Now().Format("2006/01/02 15:04:05"))
			for _, item := range msg.MonitorResponseItems {
				methodKey := fmt.Sprintf("%s.%s()", item.GetSdid(), item.GetMethod())
				logger.Blue(methodKey)
				logger.Blue(fmt.Sprintf("Total: %d, Success: %d, Fail: %d, AvgRT: %.2fus, FailRate: %.2f%%",
					item.GetTotal(), item.GetSuccess(), item.GetFail(), item.GetAvgRT(), item.GetFailRate()*100))

				allMonitorResponseItemsLock.Lock()
				if v, ok := allMonitorResponseItemsMap[methodKey]; ok {
					allMonitorResponseItemsMap[methodKey] = append(v, item)
				} else {
					allMonitorResponseItemsMap[methodKey] = []*monitorPB.MonitorResponseItem{item}
				}
				allMonitorResponseItemsLock.Unlock()
			}
		}
	},
}

func init() {
	root.Cmd.AddCommand(monitorCommand)
	monitorCommand.Flags().IntVarP(&debugPort, "port", "p", 1999, "debug port")
	monitorCommand.Flags().StringVar(&debugHost, "host", "127.0.0.1", "debug host")
	monitorCommand.Flags().IntVarP(&interval, "interval", "i", 5, "monitor interval")
}

func getAverageFloat32(input []float32) float32 {
	length := len(input)
	if length == 0 {
		return 0
	}
	sum := float32(0)
	for _, v := range input {
		sum += v
	}
	return sum / float32(length)
}

type methodSorter []string

func (m methodSorter) Len() int {
	return len(m)
}

func (m methodSorter) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m methodSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
