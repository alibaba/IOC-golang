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
