package app

import (
	"fmt"
	"restapi-bus/helper"

	"github.com/go-redis/redis"
)

func NewRedis(host string, port string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass, // no password set
		DB:       0,    // use default DB

	})
	ping := rdb.Ping()
	if ping.Err() != nil {
		helper.PanicIfError(ping.Err())
	}

	return rdb
}
