package services

import (
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

type RedisWrapper struct {
	host      string
	redisConn redis.Conn
}

func NewRedisWrapper() *RedisWrapper {
	var redisWrapper RedisWrapper

	hostEnv := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	if hostEnv != "" {
		redisWrapper.host = hostEnv
	} else {
		redisWrapper.host = "localhost"
	}

	redisConn, err := redis.Dial("tcp", fmt.Sprintf("%s:6379", redisWrapper.host))
	if err != nil {
		log.Fatal("Can't connect to redis: ", err)
	}

	redisWrapper.redisConn = redisConn

	return &redisWrapper
}

func (dw RedisWrapper) CloseConnection() {
	dw.redisConn.Close()
}
