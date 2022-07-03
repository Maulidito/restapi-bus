package main

import (
	"restapi-bus/depedency"
)

func main() {

	server := depedency.InitializedServer()

	server.Run(":8080")
}
