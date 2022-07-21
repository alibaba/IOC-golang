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

package monitor

import (
	oriCtx "context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	monitorPB "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

type mockMonitorServiceMonitorServerImpl struct {
	grpc.ServerStream
	rspList   []*monitorPB.MonitorResponse
	returnErr error
	cancel    oriCtx.CancelFunc
}

func (m *mockMonitorServiceMonitorServerImpl) Send(rsp *monitorPB.MonitorResponse) error {
	m.rspList = append(m.rspList, rsp)
	return m.returnErr
}

func (m *mockMonitorServiceMonitorServerImpl) Context() oriCtx.Context {
	ctx, cancel := oriCtx.WithCancel(oriCtx.Background())
	m.cancel = cancel
	return ctx
}

func TestMonitorServiceImpl(t *testing.T) {
	t.Run("monitor with target sdid and method success", func(t *testing.T) {
		mockInterceptorImpl := newMockInterceptor(t)

		sdid := "github.com/alibaba/ioc-golang/extension/aop/monitor/test.Struct1"
		method := "TestMethod1"
		timestamp := time.Now().UnixMicro()
		interval := int64(1)

		mockInterceptorImpl.On("Monitor", mock.MatchedBy(func(monitorCtx *context) bool {
			if monitorCtx.SDID != sdid {
				return false
			}
			if monitorCtx.MethodName != method {
				return false
			}
			return true
		})).Once().Run(func(args mock.Arguments) {
			ctx := args[0].(*context)
			go func() {
				ctx.Ch <- &monitorPB.MonitorResponse{
					MonitorResponseItems: []*monitorPB.MonitorResponseItem{
						{
							Sdid:      sdid,
							Method:    method,
							Total:     10,
							Timestamp: timestamp,
						},
					},
				}
			}()
		})

		mockInterceptorImpl.On("StopMonitor").Once()

		mockService := newMockMonitorService(mockInterceptorImpl)

		req := &monitorPB.MonitorRequest{
			Sdid:     sdid,
			Method:   method,
			Interval: interval,
		}
		serverStream := &mockMonitorServiceMonitorServerImpl{
			returnErr: nil,
			rspList:   make([]*monitorPB.MonitorResponse, 0),
		}
		go func() {
			time.AfterFunc(time.Millisecond*100, func() {
				assert.NotNil(t, serverStream.cancel)
				serverStream.cancel()
			})
		}()
		assert.Nil(t, mockService.Monitor(req, serverStream))
		assert.Equal(t, 1, len(serverStream.rspList))
		assert.Equal(t, 1, len(serverStream.rspList[0].MonitorResponseItems))
		item := serverStream.rspList[0].MonitorResponseItems[0]
		assert.Equal(t, timestamp, item.Timestamp)
		assert.Equal(t, method, item.Method)
		assert.Equal(t, sdid, item.Sdid)
		time.Sleep(time.Millisecond * 100)
	})
}
