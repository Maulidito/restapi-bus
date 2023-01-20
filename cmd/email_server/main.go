package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"restapi-bus/app"
	"restapi-bus/helper"
	"restapi-bus/models/response"
	"restapi-bus/repository"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
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
			log.Print(detailTicket)
			SendTicketEmail(&detailTicket)
			if err != nil {
				message.Reject(false)
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}

	}()

	<-stopService
}

func SendTicketEmail(detailTicket *response.DetailTicket) {
	from := mail.NewEmail("Bus Agency", "busagencyapi@gmail.com")

	subject := fmt.Sprintf("Bus Ticket %d", detailTicket.TicketId)

	htmlContent := getHtmlTemplate(detailTicket)
	to := mail.NewEmail("Bus Agency Ticket", detailTicket.Customer.Email)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Success")
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}

func getHtmlTemplate(data *response.DetailTicket) string {
	template, err := template.ParseFiles("../../views/html/email_template.html")
	helper.PanicIfError(err)
	templateBuffer := new(bytes.Buffer)

	err = template.Execute(templateBuffer, data)

	helper.PanicIfError(err)

	return templateBuffer.String()
}
