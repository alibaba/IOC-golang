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

package trace

import (
	"github.com/fatih/color"

	"github.com/alibaba/ioc-golang/aop/common"
	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
)

type traceServiceImpl struct {
	tracePB.UnimplementedTraceServiceServer
	traceInterceptor *methodTraceInterceptor
}

func (d *traceServiceImpl) Trace(req *tracePB.TraceRequest, traceServer tracePB.TraceService_TraceServer) error {
	color.Red("[Debug Server] Receive trace request %+v\n", req.String())
	defer color.Red("[Debug Server] Trace request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *tracePB.TraceResponse)
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	traceCtx := newTraceByMethodContext(sdid, method, sendCh, fieldMatcher)
	d.traceInterceptor.StartTraceByMethod(traceCtx)

	done := traceServer.Context().Done()
	if err := traceServer.Send(&tracePB.TraceResponse{
		CollectorAddress: getCollectorAddress(),
	}); err != nil {
		return err
	}

	<-done
	d.traceInterceptor.StopTraceByMethod(traceCtx)

	// todo return trace info to cli
	//case traceRsp := <-sendCh:
	//	if err := traceServer.Send(traceRsp); err != nil {
	//		return err
	//	}
	return nil
}

func newTraceGRPCService() *traceServiceImpl {
	return &traceServiceImpl{
		traceInterceptor: getTraceInterceptorSingleton(),
	}
}
