package app

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedis(host string, port string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass, // no password set
		DB:       0,    // use default DB

	})

	return rdb
}
