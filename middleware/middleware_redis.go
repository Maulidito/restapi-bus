package middleware

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"restapi-bus/models/web"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type RedisClientDb struct {
	*redis.Client
}

func (redis *RedisClientDb) MiddlewareGetDataRedis(ctx *gin.Context) {

	redisResult := redis.Get(ctx.Request.URL.String())

	res, err := redisResult.Bytes()

	if err != nil {
		fmt.Println("ERROR ", err)
		return
	}
	finalResposen := web.WebResponse{}

	byteBuffer := bytes.NewBuffer(res)

	err = gob.NewDecoder(byteBuffer).Decode(&finalResposen)

	if err != nil {
		fmt.Println("ERROR GET DATA REDIS", err, "\n", finalResposen, res)
		return
	}

	log.Println("GET DATA REDIS :", ctx.Request.URL.String())
	ctx.AbortWithStatusJSON(http.StatusOK, finalResposen)

}

func (redis *RedisClientDb) MiddlewareSetDataRedis(ctx *gin.Context) {

	timeDuration := time.Hour * 1
	go func() {
		response, ok := ctx.Get("response")

		if !ok {
			return
		}
		response1 := response.(web.WebResponse)
		byteBuffer := new(bytes.Buffer)
		if err := gob.NewEncoder(byteBuffer).Encode(&response1); err != nil {
			log.Print("ERROR ENCODE", err)
			return
		}

		log.Println("SET KEY REDIS", ctx.Request.URL.String())

		redis.Set(ctx.Request.URL.String(), byteBuffer.Bytes(), timeDuration)

	}()

	ctx.Next()

}
