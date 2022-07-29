/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,马克
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/admin"

	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=PushConsumerParam
// +ioc:autowire:constructFunc=New

type PushConsumer struct {
	rocketmq.PushConsumer
}

func (i *PushConsumer) Start() error {
	return i.PushConsumer.Start()
}

func (i *PushConsumer) Shutdown() error {
	return i.PushConsumer.Shutdown()
}

func (i *PushConsumer) Subscribe(topic string, selector consumer.MessageSelector, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	return i.PushConsumer.Subscribe(topic, selector, f)
}

func (i *PushConsumer) Unsubscribe(topic string) error {
	return i.PushConsumer.Unsubscribe(topic)
}

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=ProducerParam
// +ioc:autowire:constructFunc=New

type Producer struct {
	rocketmq.Producer
}

func (i *Producer) Start() error {
	return i.Producer.Start()
}

func (i *Producer) Shutdown() error {
	return i.Producer.Shutdown()
}

func (i *Producer) SendSync(ctx context.Context, mq ...*primitive.Message) (*primitive.SendResult, error) {
	return i.Producer.SendSync(ctx, mq...)
}

func (i *Producer) SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error), msg ...*primitive.Message) error {
	return i.Producer.SendAsync(ctx, mq, msg...)
}

func (i *Producer) SendOneWay(ctx context.Context, mq ...*primitive.Message) error {
	return i.Producer.SendOneWay(ctx, mq...)
}

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=AdminParam
// +ioc:autowire:constructFunc=New

type Admin struct {
	admin.Admin
}

func (i *Admin) CreateTopic(ctx context.Context, opts ...admin.OptionCreate) error {
	return i.Admin.CreateTopic(ctx, opts...)
}
func (i *Admin) DeleteTopic(ctx context.Context, opts ...admin.OptionDelete) error {
	return i.Admin.DeleteTopic(ctx, opts...)
}
func (i *Admin) Close() error {
	return i.Admin.Close()
}
