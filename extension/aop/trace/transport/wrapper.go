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
	"sync"

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
	collectorLock         sync.Mutex
	processor             *tJaegerClient.Process
	tJaegerClientBatchOut chan *bytes.Buffer
}

func newLocalWrappedHTTPTransport(httpTransport *transport.HTTPTransport, tJaegerClientBatchOut chan *bytes.Buffer) *localWrappedHTTPTransport {
	return &localWrappedHTTPTransport{
		httpTransport:         httpTransport,
		spans:                 make([]*tJaegerClient.Span, 0),
		tJaegerClientBatchOut: tJaegerClientBatchOut,
	}
}
func (l *localWrappedHTTPTransport) SetCollector(appName string, traceOut chan []*model.Trace, interval int) {
	c, err := newCollector(appName, interval, traceOut)
	if err != nil {
		// todo
		panic(err)
	}
	l.collectorLock.Lock()
	defer l.collectorLock.Unlock()
	l.collector = c
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
			l.collectorLock.Lock()
			if remoteCollector := l.collector; remoteCollector != nil {
				_ = remoteCollector.handle([]*tJaeger.Batch{batch})
			}
			l.collectorLock.Unlock()
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
	l.RemoveCollector()
	localWrappedHTTPTransportSingleton = nil
	return l.httpTransport.Close()
}

func (l *localWrappedHTTPTransport) RemoveCollector() {
	l.collectorLock.Lock()
	defer l.collectorLock.Unlock()
	if l.collector != nil {
		l.collector.destroy()
		l.collector = nil
	}
}

var localWrappedHTTPTransportSingleton *localWrappedHTTPTransport

func GetLocalWrappedHTTPTransportSingleton(collectorAddress string, batchBufferOut chan *bytes.Buffer) jaeger.Transport {
	if localWrappedHTTPTransportSingleton == nil {
		httpTransport := transport.NewHTTPTransport(common.GetJaegerCollectorEndpoint(collectorAddress))
		localWrappedHTTPTransportSingleton = newLocalWrappedHTTPTransport(httpTransport, batchBufferOut)
	}
	return localWrappedHTTPTransportSingleton
}

func RemoveCollector() {
	if localWrappedHTTPTransportSingleton != nil {
		localWrappedHTTPTransportSingleton.RemoveCollector()
	}
}

func SetCollector(appName string, traceOut chan []*model.Trace, interval int) {
	localWrappedHTTPTransportSingleton.SetCollector(appName, traceOut, interval)
}

func serializeThrift(obj uberThrift.TStruct) (*bytes.Buffer, error) {
	t := uberThrift.NewTMemoryBuffer()
	p := uberThrift.NewTBinaryProtocolConf(t, nil)
	if err := obj.Write(context.Background(), p); err != nil {
		return nil, err
	}
	return t.Buffer, nil
}
