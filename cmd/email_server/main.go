package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"restapi-bus/app"
	"restapi-bus/helper"
	"restapi-bus/models/response"
	"restapi-bus/repository"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {

	template := parseHtmlFile("../../views/html/email_template_v2.html", "../../views/html/email_ticket.html", "../../views/html/email_template.html")
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
			log.Print(detailTicket)
			templateRendered := renderHtmlTemplate(&detailTicket, "email_template.html", template)
			//imageHtml := htmlToImage(templateRendered)
			//log.Println("IMAGE HTML", imageHtml)
			//templateFinal := renderHtmlTemplate(&struct{ Url string }{Url: imageHtml}, "email_template_v2.html", template)
			//log.Println("Template Final \n", templateFinal)
			SendTicketEmail(&detailTicket, templateRendered)
			if err != nil {
				message.Reject(false)
			}
			err = message.Ack(true)
			helper.PanicIfError(err)
		}

	}()

	<-stopService
}

func htmlToImage(templateRendered string) string {
	dataHTML := map[string]string{
		"html": templateRendered,
	}
	templateRenderedByte, err := json.Marshal(&dataHTML)
	helper.PanicIfError(err)
	req, err := http.NewRequest("POST", "https://hcti.io/v1/image", bytes.NewReader(templateRenderedByte))

	helper.PanicIfError(err)

	req.SetBasicAuth(os.Getenv("HTMLTOIMAGE_USERID"), os.Getenv("HTMLTOIMAGE_APIKEY"))

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)

	helper.PanicIfError(err)

	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)

	helper.PanicIfError(err)

	BodyStruct := struct {
		Url string
	}{}

	if err := json.Unmarshal(bodyByte, &BodyStruct); err != nil {
		helper.PanicIfError(err)
	}

	log.Print("RESULT BODY STRING", BodyStruct.Url)
	return BodyStruct.Url

}

func SendTicketEmail(detailTicket *response.DetailTicket, templateRendered string) {
	from := mail.NewEmail("Bus Agency", "busagencyapi@gmail.com")

	subject := fmt.Sprintf("Bus Ticket %d", detailTicket.TicketId)

	to := mail.NewEmail("Bus Agency Ticket", detailTicket.Customer.Email)

	message := mail.NewSingleEmail(from, subject, to, "", templateRendered)

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

func parseHtmlFile(filePath ...string) *template.Template {
	template, err := template.ParseFiles(filePath...)
	helper.PanicIfError(err)
	return template
}

func renderHtmlTemplate(data interface{}, templateName string, template *template.Template) string {

	templateBuffer := new(bytes.Buffer)

	err := template.ExecuteTemplate(templateBuffer, templateName, data)

	helper.PanicIfError(err)

	return templateBuffer.String()
}
