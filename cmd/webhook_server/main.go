package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"restapi-bus/app"
	"restapi-bus/constant"
	"restapi-bus/external"
	"restapi-bus/helper"
	"restapi-bus/models/web"
	"restapi-bus/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var mqChannelSingleton repository.IMessageChannel

func main() {
	err := godotenv.Load("../../.env")
	helper.PanicIfError(err)
	g := app.DefaultConfigurationRouter()
	payment := external.NewPayment()
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

	connRmq := app.NewRabbitMqConn(usernameRmq, passwordRmq, hostRmq, portRmq)

	channel, err := connRmq.Channel()
	helper.PanicIfError(err)
	defer func() {
		err := channel.Close()
		helper.PanicIfError(err)
	}()

	MqChannel := repository.BindMqChannel(channel)
	mqChannelSingleton = MqChannel
	groupPayment := g.Group("/hook/payment")
	groupPayment.POST("/success", PaymentSuccessHandler)
	ngrokTunnel := ConfigNgrok()

	err = payment.SetXenditWebhookUrl("fva_paid", ngrokTunnel.URL())
	if err != nil {
		fmt.Printf("Error Set Xendit Webhook Url \n %s", err.Error())
		os.Exit(1)
	}

	g.RunListener(ngrokTunnel)

}

func ConfigNgrok() ngrok.Tunnel {
	ctx := context.Background()
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	helper.PanicIfError(err)

	return tun

}

func PaymentSuccessHandler(ctx *gin.Context) {
	webhook_verfication := os.Getenv("WEBHOOK_VERIFICATION_TOKEN")
	if webhook_verfication == "" {

		fmt.Println(fmt.Errorf("webhook verification in server is null"))

		return
	}
	webhook_token := ctx.GetHeader("x-callback-token")
	if webhook_verfication != webhook_token {
		fmt.Println(fmt.Errorf("got wrong webhook token : %s", webhook_token))

		return
	}
	paymentSuccess := web.PaymentSuccess{}
	err := ctx.ShouldBind(&paymentSuccess)
	fmt.Println(err)
	helper.PanicIfError(err)

	jsonPaymentSuccess, err := json.Marshal(paymentSuccess)
	fmt.Println(err)
	helper.PanicIfError(err)

	mqChannelSingleton.PublishToEmailServiceTopic(ctx, constant.TOPIC_PAYMENT_WEBHOOK, constant.QUEUE_PAYMENT_WEBHOOK, jsonPaymentSuccess)
	ctx.JSON(http.StatusOK, web.WebResponseNoData{Code: http.StatusOK, Status: http.StatusText(http.StatusOK)})

}
