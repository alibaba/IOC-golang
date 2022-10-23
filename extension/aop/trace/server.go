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
	"bytes"
	"sort"

	"github.com/jaegertracing/jaeger/model"

	"github.com/alibaba/ioc-golang/aop/common"
	tracePB "github.com/alibaba/ioc-golang/extension/aop/trace/api/ioc_golang/aop/trace"
	"github.com/alibaba/ioc-golang/logger"
)

const outBatchBufferAndChSize = 1000

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy:autoInjection=false

type traceServiceImpl struct {
	tracePB.UnimplementedTraceServiceServer
	TraceInterceptor traceInterceptorIOCInterface `singleton:""`
}

func (d *traceServiceImpl) Trace(req *tracePB.TraceRequest, traceServer tracePB.TraceService_TraceServer) error {
	logger.Red("[Debug Server] Receive trace request %+v\n", req.String())
	defer logger.Red("[Debug Server] Trace request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	traceCtx, _ := GetdebugServerTraceByMethodContext(&debugServerTraceByMethodContextParam{
		sdid:         sdid,
		method:       method,
		fieldMatcher: fieldMatcher,
		maxDepth:     req.MaxDepth,
		maxLength:    req.MaxLength,
	})
	d.TraceInterceptor.StartTraceByMethod(traceCtx)

	done := traceServer.Context().Done()
	if err := traceServer.Send(&tracePB.TraceResponse{
		CollectorAddress: getCollectorAddress(),
	}); err != nil {
		return err
	}

	if req.GetPushToCollectorAddress() != "" {
		// start subscribing batch buffer
		outBatchBuffer := make(chan *bytes.Buffer, outBatchBufferAndChSize)
		getGlobalTracer().subscribeBatchBuffer(outBatchBuffer)
		go func() {
			for {
				select {
				case <-done:
					getGlobalTracer().removeSubscribingBatchBuffer()
					return
				case batchBuffer := <-outBatchBuffer:
					_ = traceServer.Send(&tracePB.TraceResponse{
						ThriftSerializedSpans: batchBuffer.Bytes(),
					})
				}
			}
		}()
	}

	outTraceCh := make(chan []*model.Trace, outBatchBufferAndChSize)
	// start subscribing traces info
	getGlobalTracer().subscribeTrace(outTraceCh)
	go func() {
		for {
			select {
			case <-done:
				getGlobalTracer().removeSubscribingTrace()
				return
			case traces := <-outTraceCh:
				sortableTraces := traceSorter(traces)
				sort.Sort(sortableTraces)
				_ = traceServer.Send(&tracePB.TraceResponse{
					Traces: sortableTraces,
				})
			}
		}
	}()
	<-done
	d.TraceInterceptor.StopTraceByMethod()
	return nil
}
