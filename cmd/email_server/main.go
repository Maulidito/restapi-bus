package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"restapi-bus/app"
	"restapi-bus/cmd/email_server/service"
	"restapi-bus/helper"
	"restapi-bus/models/response"
	"restapi-bus/repository"

	"github.com/joho/godotenv"
)

func main() {

	template := service.ParseHtmlFile("../../views/html/email_template_v2.html", "../../views/html/email_ticket.html", "../../views/html/email_template.html")
	CONSUMER_NAME := "EMAIL_SERVICE"
	err := godotenv.Load("../../.env")
	helper.PanicIfError(err)

	usernameRmq := os.Getenv("USERNAME_RMQ")
	passwordRmq := os.Getenv("PASSWORD_RMQ")
	hostRmq := os.Getenv("HOST_RMQ")
	portRmq := os.Getenv("PORT_RMQ")

	conn := app.NewRabbitMqConn(usernameRmq, passwordRmq, hostRmq, portRmq)

	channel, err := conn.Channel()

	defer func() {
		err := channel.Close()
		helper.PanicIfError(err)
	}()

	helper.PanicIfError(err)

	MqChannel := repository.BindMqChannel(channel)

	messageChannel := MqChannel.ConsumeQueue(context.Background(), CONSUMER_NAME)
	helper.PanicIfError(err)

	stopService := make(chan bool)

	go func() {
		log.Print("EMAIL SERVICE STARTED")
		for message := range messageChannel {
			detailTicket := response.DetailTicket{}
			err := json.Unmarshal(message.Body, &detailTicket)
			helper.PanicIfError(err)
			templateRendered := service.RenderHtmlTemplate(&detailTicket, "email_ticket.html", template)
			service.SendTicketEmailSmtp(&detailTicket, templateRendered)
			if err != nil {
				message.Reject(false)
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}
		close(stopService)
		return

	}()

	<-stopService
}
