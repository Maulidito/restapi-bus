package repository

import (
	"context"
	"restapi-bus/helper"

	"github.com/rabbitmq/amqp091-go"
)

var MessageChannelSingleton *MessageChannel

type MessageChannel struct {
	*amqp091.Channel
}

type IMessageChannel interface {
	PublishToEmailService(ctx context.Context, queueName string, data []byte)
	PublishToEmailServiceTopic(ctx context.Context, topic string, queueName string, data []byte)
	ConsumeQueue(ctx context.Context, consumerName string, queueName string) <-chan amqp091.Delivery
	AckMessage(ctx context.Context, msg amqp091.Delivery) error
}

func BindMqChannel(channelMq *amqp091.Channel) IMessageChannel {
	if MessageChannelSingleton == nil {
		MessageChannelSingleton = &MessageChannel{channelMq}
	}
	return &MessageChannel{channelMq}
}

func (mq *MessageChannel) PublishToEmailServiceTopic(ctx context.Context, topic string, queueName string, data []byte) {
	err := mq.QueueBind(queueName, topic, "amq.topic", false, nil)
	helper.PanicIfError(err)
	mq.PublishWithContext(ctx, "amq.topic", topic, false, false, amqp091.Publishing{Body: data})
}

func (mq *MessageChannel) PublishToEmailService(ctx context.Context, queueName string, data []byte) {

	err := mq.QueueBind(queueName, "info", "amq.direct", false, nil)

	helper.PanicIfError(err)
	mq.PublishWithContext(ctx, "amq.direct", "info", false, false, amqp091.Publishing{Body: data})
}

func (mq *MessageChannel) ConsumeQueue(ctx context.Context, consumerName string, queueName string) <-chan amqp091.Delivery {
	queue, err := mq.QueueDeclare(queueName, true, false, false, true, nil)

	helper.PanicIfError(err)

	dataConsume, err := mq.Consume(queue.Name, consumerName, false, false, false, true, nil)

	helper.PanicIfError(err)

	return dataConsume
}

func (mq *MessageChannel) AckMessage(ctx context.Context, msg amqp091.Delivery) error {
	err := mq.Ack(msg.DeliveryTag, false)
	return err
}
