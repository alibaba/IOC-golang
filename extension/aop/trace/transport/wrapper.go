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

package transport

import (
	"bytes"
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/model"
	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/uber/jaeger-client-go"
	uberThrift "github.com/uber/jaeger-client-go/thrift"
	tJaegerClient "github.com/uber/jaeger-client-go/thrift-gen/jaeger"
	"github.com/uber/jaeger-client-go/transport"

	"github.com/alibaba/ioc-golang/extension/aop/trace/common"
)

type localWrappedHTTPTransport struct {
	httpTransport         *transport.HTTPTransport
	spans                 []*tJaegerClient.Span
	collector             *collector
	processor             *tJaegerClient.Process
	tJaegerClientBatchOut chan *bytes.Buffer
}

func newLocalWrappedHTTPTransport(appName string, httpTransport *transport.HTTPTransport, traceOut chan []*model.Trace, tJaegerClientBatchOut chan *bytes.Buffer, interval int) *localWrappedHTTPTransport {
	c, err := newCollector(appName, interval, traceOut)
	if err != nil {
		// todo
		panic(err)
	}
	return &localWrappedHTTPTransport{
		httpTransport:         httpTransport,
		spans:                 make([]*tJaegerClient.Span, 0),
		collector:             c,
		tJaegerClientBatchOut: tJaegerClientBatchOut,
	}
}

func (l *localWrappedHTTPTransport) Append(span *jaeger.Span) (int, error) {
	jSpan := jaeger.BuildJaegerThrift(span)
	l.spans = append(l.spans, jSpan)
	if l.processor == nil {
		l.processor = jaeger.BuildJaegerProcessThrift(span)
	}
	return l.httpTransport.Append(span)
}

func (l *localWrappedHTTPTransport) Flush() (int, error) {
	if len(l.spans) > 0 {
		// 1. store to memory collector
		tJaegerClientBatch := &tJaegerClient.Batch{
			Spans:   l.spans,
			Process: l.processor,
		}
		data, err := serializeThrift(tJaegerClientBatch)
		if err != nil {
			panic(err)
		}
		tdes := thrift.NewTDeserializer()
		batch := &tJaeger.Batch{}
		if err := tdes.Read(context.Background(), batch, data.Bytes()); err == nil {
			_ = l.collector.handle([]*tJaeger.Batch{batch})
		}
		l.spans = l.spans[:0]
		// 2. send tJaegerClientBatch bytes back to cli in order to push to 'pushAddr' if needed
		if l.tJaegerClientBatchOut != nil {
			l.tJaegerClientBatchOut <- data
		}
	}
	return l.httpTransport.Flush()
}

func (l *localWrappedHTTPTransport) Close() error {
	return l.httpTransport.Close()
}

var localWrappedHTTPTransportSingleton *localWrappedHTTPTransport

func GetLocalWrappedHTTPTransportSingleton(appName, collectorAddress string, out chan []*model.Trace, batchBufferOut chan *bytes.Buffer, interval int) jaeger.Transport {
	if localWrappedHTTPTransportSingleton == nil {
		httpTransport := transport.NewHTTPTransport(common.GetJaegerCollectorEndpoint(collectorAddress))
		localWrappedHTTPTransportSingleton = newLocalWrappedHTTPTransport(appName, httpTransport, out, batchBufferOut, interval)
	}
	return localWrappedHTTPTransportSingleton
}

func serializeThrift(obj uberThrift.TStruct) (*bytes.Buffer, error) {
	t := uberThrift.NewTMemoryBuffer()
	p := uberThrift.NewTBinaryProtocolConf(t, nil)
	if err := obj.Write(context.Background(), p); err != nil {
		return nil, err
	}
	return t.Buffer, nil
}
