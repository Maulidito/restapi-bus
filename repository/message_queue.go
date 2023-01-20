package repository

import (
	"context"
	"restapi-bus/helper"

	"github.com/rabbitmq/amqp091-go"
)

type MessageChannel struct {
	*amqp091.Channel
}

type IMessageChannel interface {
	PublishToEmailService(ctx context.Context, data *amqp091.Publishing)
	ConsumeQueue(ctx context.Context, consumerName string) <-chan amqp091.Delivery
}

func BindMqChannel(channelMq *amqp091.Channel) IMessageChannel {

	return &MessageChannel{channelMq}
}

func (mq *MessageChannel) PublishToEmailService(ctx context.Context, data *amqp091.Publishing) {
	queue := initQueue(mq)
	err := mq.QueueBind(queue.Name, "info", "amq.direct", false, nil)
	helper.PanicIfError(err)
	mq.PublishWithContext(ctx, "amq.direct", "info", false, false, *data)
}

func (mq *MessageChannel) ConsumeQueue(ctx context.Context, consumerName string) <-chan amqp091.Delivery {
	queue := initQueue(mq)
	dataConsume, err := mq.Consume(queue.Name, consumerName, false, false, false, true, nil)

	helper.PanicIfError(err)

	return dataConsume
}

func initQueue(mq *MessageChannel) amqp091.Queue {
	queue, err := mq.QueueDeclare("busTicket", false, true, false, true, nil)
	helper.PanicIfError(err)

	return queue

}
