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

package server

import (
	"context"
	"sort"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/debug"

	debugCommon "github.com/alibaba/ioc-golang/debug/common"
	"github.com/alibaba/ioc-golang/debug/interceptor"
	"github.com/alibaba/ioc-golang/debug/interceptor/common"
	"github.com/alibaba/ioc-golang/debug/interceptor/trace"
	"github.com/alibaba/ioc-golang/debug/interceptor/watch"
)

type DebugServerImpl struct {
	traceInterceptor        *trace.Interceptor
	watchInterceptor        *watch.Interceptor
	allInterfaceMetadataMap map[string]*debugCommon.StructMetadata
	debug.UnimplementedDebugServiceServer
}

func (d *DebugServerImpl) ListServices(ctx context.Context, empty *emptypb.Empty) (*debug.ListServiceResponse, error) {
	structsMetadatas := make(interceptor.MetadataSorter, 0)
	for key, v := range d.allInterfaceMetadataMap {
		methods := make([]string, 0)
		for key := range v.MethodMetadata {
			methods = append(methods, key)
		}

		structsMetadatas = append(structsMetadatas, &debug.ServiceMetadata{
			Methods:            methods,
			InterfaceName:      key,
			ImplementationName: key,
		})
	}
	sort.Sort(structsMetadatas)

	return &debug.ListServiceResponse{
		ServiceMetadata: structsMetadatas,
	}, nil
}

func (d *DebugServerImpl) Watch(req *debug.WatchRequest, watchSever debug.DebugService_WatchServer) error {
	color.Red("[Debug Server] Receive watch request %+v\n", req.String())
	defer color.Red("[Debug Server] Watch request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *debug.WatchResponse)
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	watchCtx := &watch.Context{
		SDID:         sdid,
		MethodName:   method,
		Ch:           sendCh,
		FieldMatcher: fieldMatcher,
	}
	d.watchInterceptor.Watch(watchCtx)

	done := watchSever.Context().Done()
	for {
		select {
		case <-done:
			// watch stop
			d.watchInterceptor.UnWatch(watchCtx)
			return nil
		case watchRsp := <-sendCh:
			if err := watchSever.Send(watchRsp); err != nil {
				return err
			}
		}
	}
}

func (d *DebugServerImpl) Trace(req *debug.TraceRequest, traceServer debug.DebugService_TraceServer) error {
	color.Red("[Debug Server] Receive trace request %+v\n", req.String())
	defer color.Red("[Debug Server] Trace request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *debug.TraceResponse)
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	traceCtx := &trace.Context{
		SDID:         sdid,
		MethodName:   method,
		Ch:           sendCh,
		FieldMatcher: fieldMatcher,
	}
	d.traceInterceptor.Trace(traceCtx)

	done := traceServer.Context().Done()
	if err := traceServer.Send(&debug.TraceResponse{
		CollectorAddress: trace.GetCollectorAddress(),
	}); err != nil {
		return err
	}

	<-done
	d.traceInterceptor.UnTrace(traceCtx)

	// todo return trace info to cli
	//case traceRsp := <-sendCh:
	//	if err := traceServer.Send(traceRsp); err != nil {
	//		return err
	//	}
	return nil
}
