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
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/extension/aop/trace/goroutine_trace"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=initRPCInterceptor

type rpcInterceptor struct {
	tracingGRIDMap   map[int64]struct{}
	TraceInterceptor goroutine_trace.GoRoutineTraceInterceptorIOCInterface `singleton:""`
}

func initRPCInterceptor(r *rpcInterceptor) (*rpcInterceptor, error) {
	r.tracingGRIDMap = make(map[int64]struct{})
	return r, nil
}

func (r *rpcInterceptor) BeforeServerInvoke(c *gin.Context) error {
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	clientContext, err := getGlobalTracer().getRawTracer().Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		requestPath := c.Request.URL.Path
		splitedPath := strings.Split(requestPath, "/")
		if len(splitedPath) <= 2 {
			return fmt.Errorf("invalid request path %s", requestPath)
		}
		method := splitedPath[len(splitedPath)-1]

		// create facade ctx
		facadeCtx, err := GettraceGoRoutineInterceptorFacadeCtx(&traceGoRoutineInterceptorFacadeCtxParam{
			trace:             newTraceWithClientSpanContext(method, clientContext),
			clientSpanContext: clientContext,
		})
		if err != nil {
			log.Printf("rpc trace Interceptor GettraceGoRoutineInterceptorFacadeCtx failed with erorr = %s\n", err.Error())
			return err
		}

		traceByGrContext, err := goroutine_trace.GetGoRoutineTracingContext(&goroutine_trace.GoRoutineTracingContextParams{
			/**
			FIXME: now we just put short method name as full name, this would make TraceInterceptor never meet
			'current span' and never jump out of tracing, so we call r.TraceInterceptor.DeleteCurrentGRTracingContext()
			in AfterServerInvoke to force stop the tracing, but that's not graceful

			Now, EntranceMethodFullName can be any string except empty, we just want 'current span' never match in TraceInterceptor
			*/
			EntranceMethodFullName: method,
			FacadeCtx:              facadeCtx,
		})
		if err != nil {
			return err
		}
		r.TraceInterceptor.AddCurrentGRTracingContext(traceByGrContext)
		r.tracingGRIDMap[goid.Get()] = struct{}{}
	}
	return nil
}

func (r *rpcInterceptor) AfterServerInvoke(ctx *gin.Context) error {
	grID := goid.Get()
	if _, ok := r.tracingGRIDMap[grID]; ok {
		// force stop tracing as the rpc is finished
		// FIXME: gr tracing type is ignored, this would cause existing debug log tracing failed.
		r.TraceInterceptor.DeleteCurrentGRTracingContext()
		delete(r.tracingGRIDMap, grID)
	}
	return nil
}

func (r *rpcInterceptor) BeforeClientInvoke(req *http.Request) error {
	// inject tracing context if necessary
	if currentGRTracingCtx := r.TraceInterceptor.GetCurrentGRTracingContext(traceGoRoutineInterceptorFacadeCtxType); currentGRTracingCtx != nil {
		if currentSpan := currentGRTracingCtx.GetFacadeCtx(); currentSpan != nil {
			// current rpc invocation is in tracing link
			carrier := opentracing.HTTPHeadersCarrier(req.Header)
			_ = getGlobalTracer().getRawTracer().Inject(currentSpan.(*traceGoRoutineInterceptorFacadeCtx).clientSpanContext, opentracing.HTTPHeaders, carrier)
		}
	}
	return nil
}

func (r *rpcInterceptor) AfterClientInvoke(response *http.Response) error {
	return nil
}
