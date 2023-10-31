package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"restapi-bus/app"
	croncustom "restapi-bus/cron_custom"
	"restapi-bus/depedency"
	"restapi-bus/external"
	"restapi-bus/helper"
	"restapi-bus/middleware"
	"restapi-bus/models/response"
	"restapi-bus/models/web"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	helper.PanicIfError(err)

	gob.Register(web.WebResponse{})
	gob.Register(response.Agency{})
	gob.Register(response.Bus{})
	gob.Register(response.Ticket{})
	gob.Register(response.Driver{})
	gob.Register(response.Customer{})
	gob.Register(response.Schedule{})

	port := os.Getenv("PORT")
	usernameDb := os.Getenv("USERNAME_DB")
	portDb := os.Getenv("PORT_DB")
	hostDb := os.Getenv("HOST_DB")
	passDb := os.Getenv("PASSWORD_DB")
	nameDb := os.Getenv("NAME_DB")
	hostRdb := os.Getenv("HOST_RDB")
	passRdb := os.Getenv("PASSWORD_RDB")
	portRdb := os.Getenv("PORT_RDB")
	usernameRmq := os.Getenv("USERNAME_RMQ")
	passwordRmq := os.Getenv("PASSWORD_RMQ")
	hostRmq := os.Getenv("HOST_RMQ")
	portRmq := os.Getenv("PORT_RMQ")
	var hostEnv string
	flag.StringVar(&hostEnv, "hostenv", "local", "environment host every db")
	flag.Parse()
	if hostEnv == "local" {
		hostDb = "localhost"
		hostRdb = "localhost"
		hostRmq = "localhost"

	}
	db := app.NewDatabase(usernameDb, passDb, nameDb, hostDb, portDb)

	rdb := app.NewRedis(hostRdb, portRdb, passRdb)
	middlewareRedis := middleware.RedisClientDb{Client: rdb}
	Rabbitmq, err := app.NewRabbitMqConn(usernameRmq, passwordRmq, hostRmq, portRmq).Channel()
	helper.PanicIfError(err)
	apiPayment := external.NewPayment()
	cronJob := croncustom.NewCronJob()
	server := depedency.InitializedServer(&middlewareRedis, Rabbitmq, apiPayment, cronJob, app.NewDbManager(db))
	fmt.Println("SERVER RUNNING ON PORT ", port)
	server.Run(":" + port)
}
