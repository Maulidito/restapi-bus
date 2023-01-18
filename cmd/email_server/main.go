package main

import (
	"encoding/json"
	"log"
	"os"
	"restapi-bus/app"
	"restapi-bus/helper"
	"restapi-bus/models/response"

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

	queue, err := channel.QueueDeclare("busTicket", false, true, false, true, nil)

	helper.PanicIfError(err)
	messageChannel, err := channel.Consume(queue.Name, CONSUMER_NAME, false, false, true, true, nil)
	helper.PanicIfError(err)

	stopService := make(chan bool)

	go func() {
		log.Print("EMAIL SERVICE STARTED")
		for message := range messageChannel {
			detailTicket := response.DetailTicket{}
			err := json.Unmarshal(message.Body, &detailTicket)
			log.Print(detailTicket)
			if err != nil {
				//do something
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}

	}()

	<-stopService
}

func SendEmail() {
	from := mail.NewEmail("Bus Agency", "busagencyapi@gmail.com")

	subject := "test"
	plainTextContent := "plain text content"

	htmlContent := "<h1> HTML CONTENT </h1>"
	to := mail.NewEmail("Bus Agency Ticket", "maudana111@gmail.com")

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

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
