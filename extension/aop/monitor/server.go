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
	"time"

	"github.com/fatih/color"

	monitorPB "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

type monitorService struct {
	monitorPB.UnimplementedMonitorServiceServer
	monitorInterceptor interceptor
}

func (w *monitorService) Monitor(req *monitorPB.MonitorRequest, svr monitorPB.MonitorService_MonitorServer) error {
	color.Red("[Debug Server] Receive monitor request %s\n", req.String())
	defer color.Red("[Debug Server] Monitor %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *monitorPB.MonitorResponse)
	interval := 5
	if reqInterval := req.GetInterval(); reqInterval > 0 {
		interval = int(reqInterval)
	}

	monitorCtx := newContext(sdid, method, sendCh, time.Duration(interval)*time.Second)
	w.monitorInterceptor.Monitor(monitorCtx)

	done := svr.Context().Done()
	for {
		select {
		case <-done:
			// monitor stop
			w.monitorInterceptor.StopMonitor()
			return nil
		case monitorRsp := <-sendCh:
			if err := svr.Send(monitorRsp); err != nil {
				return err
			}
		}
	}
}

func newMonitorService() *monitorService {
	return &monitorService{
		monitorInterceptor: getMonitorInterceptorSingleton(),
	}
}

func newMockMonitorService(mockInterceptor interceptor) *monitorService {
	return &monitorService{
		monitorInterceptor: mockInterceptor,
	}
}
