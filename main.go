package main

import (
	"restapi-bus/depedency"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		server := depedency.InitializedServer()
		server.Run(":8080")
		wg.Done()
	}()

	wg.Wait()

}
