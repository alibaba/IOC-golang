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
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

var tracer opentracing.Tracer
var collectorAddress = "127.0.0.1:14268"
var appName = "ioc-golang-application"

func GetGlobalTracer() opentracing.Tracer {
	if tracer == nil {
		tracer, _ = newJaegerTracer(appName, collectorAddress)
	}
	return tracer
}

func GetCollectorAddress() string {
	return collectorAddress
}

func SetCollectorAddress(addr string) {
	collectorAddress = addr
}

func SetAppName(name string) {
	appName = name
}

func newJaegerTracer(service string, collectorAddress string) (opentracing.Tracer, io.Closer) {
	return jaeger.NewTracer(
		service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport.NewHTTPTransport(fmt.Sprintf("http://%s/api/traces?format=jaeger.thrift", collectorAddress))),
	)
}
