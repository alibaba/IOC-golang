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

package rocketmq

import (
	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type ProducerParam struct {
	NameServer primitive.NamesrvAddr
	GroupName  string
	Retry      int
	Options    []producer.Option
}

func (c *ProducerParam) New(impl *Producer) (*Producer, error) {
	opts := c.Options
	if len(c.NameServer) > 0 {
		opts = append(opts, producer.WithNameServer(c.NameServer))
	}
	if c.GroupName != "" {
		opts = append(opts, producer.WithGroupName(c.GroupName))
	}
	if c.Retry > 0 {
		opts = append(opts, producer.WithRetry(c.Retry))
	}

	newProducer, err := rocketmq.NewProducer(opts...)
	if err != nil {
		return nil, err
	}
	impl.Producer = newProducer
	return impl, nil
}

type PushConsumerParam struct {
	NameServer primitive.NamesrvAddr
	GroupName  string
	Retry      int
	Options    []consumer.Option
}

func (c *PushConsumerParam) New(impl *PushConsumer) (*PushConsumer, error) {
	opts := c.Options
	if len(c.NameServer) > 0 {
		opts = append(opts, consumer.WithNameServer(c.NameServer))
	}
	if c.GroupName != "" {
		opts = append(opts, consumer.WithGroupName(c.GroupName))
	}
	if c.Retry > 0 {
		opts = append(opts, consumer.WithRetry(c.Retry))
	}

	newPushConsumer, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		return nil, err
	}
	impl.PushConsumer = newPushConsumer
	return impl, nil
}

type AdminParam struct {
	NameServer primitive.NamesrvAddr
	Options    []admin.AdminOption
}

func (c *AdminParam) New(impl *Admin) (*Admin, error) {
	opts := c.Options
	if len(c.NameServer) > 0 {
		opts = append(opts, admin.WithResolver(primitive.NewPassthroughResolver(c.NameServer)))
	}

	newAdmin, err := admin.NewAdmin(opts...)
	if err != nil {
		return nil, err
	}
	impl.Admin = newAdmin
	return impl, nil
}
