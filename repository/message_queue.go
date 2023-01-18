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
}

func BindMqChannel(channelMq *amqp091.Channel) IMessageChannel {

	return &MessageChannel{channelMq}
}

func (mq *MessageChannel) PublishToEmailService(ctx context.Context, data *amqp091.Publishing) {

	queue, err := mq.QueueDeclare("busTicket", false, true, false, true, nil)
	helper.PanicIfError(err)
	err = mq.QueueBind(queue.Name, "info", "amq.direct", false, nil)
	helper.PanicIfError(err)
	mq.PublishWithContext(ctx, "amq.direct", "info", false, false, *data)
}
