package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"restapi-bus/helper"
	"restapi-bus/models/response"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

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

func SendTicketEmailSmtp(detailTicket *response.DetailTicket, templateRendered string) {

	app_password := os.Getenv("APP_PASSWORD_GMAIL")
	host_smtp := os.Getenv("SMTP_MAIL_SERVER")
	port_stmp := os.Getenv("SMTP_MAIL_PORT")

	to := detailTicket.Customer.Email

	auth := smtp.PlainAuth("", "maudana111restapibus@gmail.com", app_password, host_smtp)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host_smtp,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host_smtp, port_stmp), tlsconfig)

	helper.PanicIfError(err)

	c, _ := smtp.NewClient(conn, host_smtp)
	defer c.Close()
	err = c.Auth(auth)
	helper.PanicIfError(err)
	err = c.Mail("maudana111restapibus@gmail.com")
	helper.PanicIfError(err)
	err = c.Rcpt(to)
	helper.PanicIfError(err)
	w, err := c.Data()
	helper.PanicIfError(err)

	subject := fmt.Sprintf("Subject: Bus Ticket %d", detailTicket.TicketId)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	message := fmt.Sprintf("To: %s\r\n"+
		"%s\r\n"+
		"%s\r\n"+
		"%s\r\n", detailTicket.Customer.Email, subject, mime, templateRendered)

	_, err = w.Write([]byte(message))

	helper.PanicIfError(err)

	err = w.Close()

	helper.PanicIfError(err)

	log.Printf("Success %s\n", subject)
}

func ParseHtmlFile(filePath ...string) *template.Template {
	template, err := template.ParseFiles(filePath...)
	helper.PanicIfError(err)
	return template
}

func RenderHtmlTemplate(data interface{}, templateName string, template *template.Template) string {

	templateBuffer := new(bytes.Buffer)

	err := template.ExecuteTemplate(templateBuffer, templateName, data)

	helper.PanicIfError(err)

	return templateBuffer.String()
}
