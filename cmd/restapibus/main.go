package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"restapi-bus/app"
	"restapi-bus/depedency"
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
	db := app.NewDatabase(usernameDb, passDb, nameDb, hostDb)

	rdb := app.NewRedis(hostRdb, portRdb, passRdb)
	middlewareRedis := middleware.RedisClientDb{Client: rdb}
	Rabbitmq, err := app.NewRabbitMqConn(usernameRmq, passwordRmq, hostRmq, portRmq).Channel()
	helper.PanicIfError(err)

	server := depedency.InitializedServer(db, &middlewareRedis, Rabbitmq)
	fmt.Println("SERVER RUNNING ON PORT ", port)
	server.Run(":" + port)

}
