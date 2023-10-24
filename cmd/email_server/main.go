package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"restapi-bus/app"
	"restapi-bus/cmd/email_server/service"
	"restapi-bus/constant"
	"restapi-bus/helper"
	"restapi-bus/models/response"
	"restapi-bus/repository"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	//TODO CONSUME WEBHOOK
	template := service.ParseHtmlFile(
		"../../views/html/email_template_v2.html",
		"../../views/html/email_ticket.html",
		"../../views/html/email_template.html",
		"../../views/html/email_order_payment.html",
	)

	err := godotenv.Load("../../.env")
	helper.PanicIfError(err)
	usernameRmq := os.Getenv("USERNAME_RMQ")
	passwordRmq := os.Getenv("PASSWORD_RMQ")
	hostRmq := os.Getenv("HOST_RMQ")
	portRmq := os.Getenv("PORT_RMQ")
	var hostEnv string
	flag.StringVar(&hostEnv, "hostenv", "local", "environment host every db")
	flag.Parse()
	if hostEnv == "local" {
		hostRmq = "localhost"
	}

	conn := app.NewRabbitMqConn(usernameRmq, passwordRmq, hostRmq, portRmq)

	channel, err := conn.Channel()

	defer func() {
		err := channel.Close()
		helper.PanicIfError(err)
	}()

	helper.PanicIfError(err)

	MqChannel := repository.BindMqChannel(channel)

	wg := sync.WaitGroup{}
	wg.Add(2)
	messageChannelTicket := MqChannel.ConsumeQueue(context.Background(), constant.CONSUMER_NAME_TICKET, constant.QUEUE_TICKET)

	messageChannelPayment := MqChannel.ConsumeQueue(context.Background(), constant.CONSUMER_NAME_PAYMENT, constant.QUEUE_PAYMENT)

	go func() {
		defer wg.Done()
		log.Print("EMAIL SERVICE STARTED TICKET")
		for message := range messageChannelTicket {
			detailTicket := response.DetailTicket{}
			err := json.Unmarshal(message.Body, &detailTicket)
			helper.PanicIfError(err)

			templateRendered := service.RenderHtmlTemplate(&detailTicket, "email_ticket.html", template)
			service.SendDataEmailSmtp(
				templateRendered,
				fmt.Sprintf("Bus Ticket %d", detailTicket.TicketId),
				detailTicket.Customer.Email)
			if err != nil {
				message.Reject(false)
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}

	}()

	go func() {
		defer wg.Done()
		log.Print("EMAIL SERVICE STARTED PAYMENT")
		for message := range messageChannelPayment {
			TicketOrder := response.TicketOrder{}
			err := json.Unmarshal(message.Body, &TicketOrder)
			helper.PanicIfError(err)

			templateRendered := service.RenderHtmlTemplate(&TicketOrder, "email_order_payment.html", template)
			service.SendDataEmailSmtp(
				templateRendered,
				"Payment Order",
				TicketOrder.Customer.Email)
			if err != nil {
				message.Reject(false)
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}

	}()

	wg.Wait()
}
