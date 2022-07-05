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
	"fmt"
	"net/http"
	"strings"

	"github.com/petermattis/goid"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/debug/interceptor"
)

type rpcInterceptor struct {
	tracingGRIDMap map[int64]struct{}
}

func (r *rpcInterceptor) BeforeServerInvoke(c *gin.Context) error {
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	clientContext, err := GetGlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		requestPath := c.Request.URL.Path
		splitedPath := strings.Split(requestPath, "/")
		if len(splitedPath) <= 2 {
			return fmt.Errorf("invalid request path %s", requestPath)
		}
		method := splitedPath[len(splitedPath)-1]

		traceByGrContext := newGoRoutineTracingContextWithClientSpan(method, clientContext)
		GetTraceInterceptor().TraceCurrentGR(traceByGrContext)
		r.tracingGRIDMap[goid.Get()] = struct{}{}
	}
	return nil
}

func (r *rpcInterceptor) AfterServerInvoke(_ *gin.Context) error {
	grID := goid.Get()
	if _, ok := r.tracingGRIDMap[grID]; ok {
		GetTraceInterceptor().StopTraceCurrentGR()
		delete(r.tracingGRIDMap, grID)
	}
	return nil
}

func (r *rpcInterceptor) BeforeClientInvoke(req *http.Request) error {
	// inject tracing context if necessary
	if currentSpan := GetTraceInterceptor().GetCurrentSpan(); currentSpan != nil {
		// current rpc invocation is in tracing link
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		_ = GetGlobalTracer().Inject(currentSpan.Context(), opentracing.HTTPHeaders, carrier)
	}
	return nil
}

func (r *rpcInterceptor) AfterClientInvoke(response *http.Response) error {
	return nil
}

func init() {
	interceptor.RegisterRPCInterceptor(&rpcInterceptor{
		tracingGRIDMap: map[int64]struct{}{},
	})
}
