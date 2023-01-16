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
	err := godotenv.Load()
	helper.PanicIfError(err)

	gob.Register(web.WebResponse{})
	gob.Register(response.Agency{})
	gob.Register(response.Bus{})

	port := os.Getenv("PORT")
	usernameDb := os.Getenv("USERNAME_DB")
	hostDb := os.Getenv("HOST_DB")
	passDb := os.Getenv("PASSWORD_DB")
	nameDb := os.Getenv("NAME_DB")
	hostRdb := os.Getenv("HOST_RDB")
	passRdb := os.Getenv("PASSWORD_RDB")
	portRdb := os.Getenv("PORT_RDB")
	db := app.NewDatabase(usernameDb, passDb, nameDb, hostDb)

	rdb := app.NewRedis(hostRdb, portRdb, passRdb)
	middlewareRedis := middleware.RedisClientDb{Client: rdb}

	server := depedency.InitializedServer(db, &middlewareRedis)
	fmt.Println("SERVER RUNNING ON PORT ", port)
	server.Run(":" + port)

}
