package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	rocketmq "github.com/cinience/go_rocketmq"

	"github.com/alibaba/IOC-Golang/autowire/normal"
)

const SDID = "RocketMQClient-Impl"

type RocketMQClient interface {
	Subscribe(topic string, selector consumer.MessageSelector, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error
	Unsubscribe(topic string) error

	SendSync(ctx context.Context, mq ...*primitive.Message) (*primitive.SendResult, error)
	SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error), msg ...*primitive.Message) error
	SendOneWay(ctx context.Context, mq ...*primitive.Message) error
}

// +ioc:autowire=true
// +ioc:autowire:interface=RocketMQClient
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New

type Impl struct {
	rocketmq.PushConsumer
	rocketmq.Producer
}

var _ RocketMQClient = &Impl{}

func GetRocketMQClient(config *Config) (RocketMQClient, error) {
	rmqClientImpl, err := normal.GetImpl(SDID, config)
	if err != nil {
		return nil, err
	}
	return rmqClientImpl.(RocketMQClient), nil
}
