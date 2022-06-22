package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	rocketmq "github.com/cinience/go_rocketmq"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=New

type Impl struct {
	rocketmq.PushConsumer
	rocketmq.Producer
}

func (i *Impl) Subscribe(topic string, selector consumer.MessageSelector, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	return i.PushConsumer.Subscribe(topic, selector, f)
}

func (i *Impl) Unsubscribe(topic string) error {
	return i.PushConsumer.Unsubscribe(topic)
}

func (i *Impl) SendSync(ctx context.Context, mq ...*primitive.Message) (*primitive.SendResult, error) {
	return i.Producer.SendSync(ctx, mq...)
}

func (i *Impl) SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error), msg ...*primitive.Message) error {
	return i.Producer.SendAsync(ctx, mq, msg...)
}

func (i *Impl) SendOneWay(ctx context.Context, mq ...*primitive.Message) error {
	return i.Producer.SendOneWay(ctx, mq...)
}
