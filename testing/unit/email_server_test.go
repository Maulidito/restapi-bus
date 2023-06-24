package testing

import (
	"os"
	"restapi-bus/cmd/email_server/service"
	"restapi-bus/models/response"
	"testing"
	"text/template"

	"github.com/joho/godotenv"
)

var detailTicket response.DetailTicket
var templateText *template.Template

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	templateText = service.ParseHtmlFile("../../views/html/email_template_v2.html", "../../views/html/email_ticket.html", "../../views/html/email_template.html")

	if err != nil {
		os.Exit(400)
	}
	detailTicket = response.DetailTicket{TicketId: 12, Schedule: response.DetailSchedule{ScheduleId: 1, FromAgency: response.Agency{Place: "DEPOK"}}, Customer: response.Customer{Name: "helo", Email: "maudana111@gmail.com"}}
	code := m.Run()
	os.Exit(code)
}

func TestSendEmailSmtp(t *testing.T) {
	// t.Skip()
	templateRendered := service.RenderHtmlTemplate(&detailTicket, "email_ticket.html", templateText)
	service.SendTicketEmailSmtp(&detailTicket, templateRendered)

}

// func TestSendEmailGomail(t *testing.T) {
// 	t.Skip()
// 	service.SendTicketEmailGomail(&detailTicket, "")
// }
