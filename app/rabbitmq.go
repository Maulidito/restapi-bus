package app

import (
	"fmt"
	"restapi-bus/helper"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMqConn(username string, password string, host string, port string) *amqp091.Connection {
	println(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))

	helper.PanicIfError(err)

	return conn
}
