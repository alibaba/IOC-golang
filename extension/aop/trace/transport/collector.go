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
	"context"
	"time"

	"github.com/alibaba/ioc-golang/logger"

	"github.com/jaegertracing/jaeger/cmd/collector/app"
	"github.com/jaegertracing/jaeger/cmd/collector/app/flags"
	"github.com/jaegertracing/jaeger/cmd/collector/app/handler"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/metrics"
	"github.com/jaegertracing/jaeger/plugin/storage"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const memoryStorageType = "memory"
const memoryMaxTraces = 1000
const dynamicQueueMaxMemory = 1 * 1024 * 1024 // 1MB
const queueSize = 1000

type collector struct {
	spanHandlers   *app.SpanHandlers
	spanReader     spanstore.Reader
	stopCh         chan struct{}
	out            chan []*model.Trace
	appName        string
	storageFactory *storage.Factory
}

var storageFactory *storage.Factory
var spanHandler *app.SpanHandlers
var spanReader spanstore.Reader

func initStorageFactory() error {
	if storageFactory != nil {
		return nil
	}
	logger := zap.NewNop()
	v := viper.New()
	v.Set("memory.max-traces", memoryMaxTraces)

	newStorageFactory, err := storage.NewFactory(storage.FactoryConfig{
		SpanWriterTypes:         []string{memoryStorageType},
		SpanReaderType:          memoryStorageType,
		DependenciesStorageType: memoryStorageType,
	})
	if err != nil {
		return err
	}

	newStorageFactory.InitFromViper(v, logger)
	if err := newStorageFactory.Initialize(metrics.NullFactory, logger); err != nil {
		return err
	}
	spanWriter, err := newStorageFactory.CreateSpanWriter()
	if err != nil {
		return err
	}

	collectorOpts, err := new(flags.CollectorOptions).InitFromViper(v, zap.NewNop())
	if err != nil {
		return err
	}
	collectorOpts.DynQueueSizeMemory = dynamicQueueMaxMemory
	collectorOpts.QueueSize = queueSize

	handlerBuilder := &app.SpanHandlerBuilder{
		SpanWriter:    spanWriter,
		CollectorOpts: collectorOpts,
		Logger:        logger,
	}

	spanProcessor := handlerBuilder.BuildSpanProcessor()
	newSpanHandlers := handlerBuilder.BuildHandlers(spanProcessor)

	newSpanReader, err := newStorageFactory.CreateSpanReader()
	if err != nil {
		return err
	}
	storageFactory = newStorageFactory
	spanHandler = newSpanHandlers
	spanReader = newSpanReader
	return nil
}

func newCollector(appName string, interval int, out chan []*model.Trace) (*collector, error) {
	if err := initStorageFactory(); err != nil {
		return nil, err
	}
	newCreatedCollector := &collector{
		spanHandlers:   spanHandler,
		spanReader:     spanReader,
		stopCh:         make(chan struct{}),
		out:            out,
		appName:        appName,
		storageFactory: storageFactory,
	}
	go newCreatedCollector.runReadLoop(time.Second * time.Duration(interval))

	return newCreatedCollector, nil
}

func (c *collector) runReadLoop(period time.Duration) {
	ticker := time.NewTicker(period)
	lastTime := time.Now()
	for {
		select {
		case <-ticker.C:
			params := &spanstore.TraceQueryParameters{
				StartTimeMin: lastTime,
				ServiceName:  c.appName,
			}
			traces, err := c.spanReader.FindTraces(context.Background(), params)
			if err != nil {
				logger.Cyan("jaeger local collector read err = ", err)
				continue
			}
			lastTime = time.Now()
			if len(traces) > 0 {
				c.out <- traces
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *collector) destroy() {
	close(c.stopCh)
}

func (c *collector) handle(batches []*tJaeger.Batch) error {
	_, err := c.spanHandlers.JaegerBatchesHandler.SubmitBatches(batches, handler.SubmitBatchOptions{})
	return err
}
