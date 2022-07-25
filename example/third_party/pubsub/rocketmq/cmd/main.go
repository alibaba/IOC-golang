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

package main

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/extension/pubsub/rocketmq"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init
// +ioc:autowire:alias=AppAlias

type App struct {
	createByAPIRocketmqAdmin        rocketmq.AdminIOCInterface
	createByAPIRocketmqProducer     rocketmq.ProducerIOCInterface
	createByAPIRocketmqPushConsumer rocketmq.PushConsumerIOCInterface

	Admin        rocketmq.AdminIOCInterface        `singleton:""`
	Producer     rocketmq.ProducerIOCInterface     `singleton:""`
	PushConsumer rocketmq.PushConsumerIOCInterface `singleton:""`
}

func (a *App) Run() {
	listenAndPushMessage(a.Admin, a.Producer, a.PushConsumer)
	listenAndPushMessage(a.createByAPIRocketmqAdmin, a.createByAPIRocketmqProducer, a.createByAPIRocketmqPushConsumer)
}

type Param struct {
	NamingServer string
	GroupName    string
}

func (p *Param) Init(a *App) (*App, error) {
	createByAPIRocketmqProducer, err := rocketmq.GetProducerIOCInterface(&rocketmq.ProducerParam{
		GroupName:  p.GroupName,
		NameServer: primitive.NamesrvAddr{p.NamingServer},
	})
	if err != nil {
		panic(err)
	}
	a.createByAPIRocketmqProducer = createByAPIRocketmqProducer

	createByAPIRocketmqPushConsumer, err := rocketmq.GetPushConsumerIOCInterface(&rocketmq.PushConsumerParam{
		GroupName:  p.GroupName,
		NameServer: primitive.NamesrvAddr{p.NamingServer},
	})
	if err != nil {
		panic(err)
	}
	a.createByAPIRocketmqPushConsumer = createByAPIRocketmqPushConsumer

	createByAPIRocketmqAdmin, err := rocketmq.GetAdminIOCInterface(&rocketmq.AdminParam{
		NameServer: primitive.NamesrvAddr{p.NamingServer},
	})
	if err != nil {
		panic(err)
	}
	a.createByAPIRocketmqAdmin = createByAPIRocketmqAdmin
	return a, nil
}

func listenAndPushMessage(rmqAdmin rocketmq.AdminIOCInterface, rmqProducer rocketmq.ProducerIOCInterface, rmqPushConsumer rocketmq.PushConsumerIOCInterface) {
	// 1. create topic
	if err := rmqAdmin.CreateTopic(context.Background(),
		admin.WithTopicCreate("mytopic"),
		admin.WithBrokerAddrCreate("127.0.0.1:10911")); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)

	// FIXME: make following logic testable
	/*
		// 2. subscribe
		if err := rmqPushConsumer.Subscribe("mytopic", consumer.MessageSelector{}, func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			fmt.Printf("subscribe callback: %v \n", msgs)
			return consumer.ConsumeSuccess, nil
		}); err != nil {
			panic(err)
		}
		if err := rmqPushConsumer.Start(); err != nil {
			panic(err)
		}
		// 3. send
		if err := rmqProducer.Start(); err != nil {
			panic(err)
		}
		_, err := rmqProducer.SendSync(context.Background(), &primitive.Message{
			Topic: "mytopic",
			Body:  []byte("hello"),
		})
		if err != nil {
			panic(err)
		}

		// 4. delete topic
		if err := rmqAdmin.DeleteTopic(context.Background(),
			admin.WithTopicDelete("mytopic"),
			admin.WithBrokerAddrDelete("127.0.0.1:10911")); err != nil {
			panic(err)
		}

	*/
}

func main() {
	if err := ioc.Load(
		config.WithSearchPath("../conf")); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton(&Param{
		NamingServer: "127.0.0.1:9876",
		GroupName:    "default",
	})
	if err != nil {
		panic(err)
	}

	app.Run()
}
