package model

import "log"
import "github.com/go-redis/redis"

var redisdb *redis.Client

func Del() {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisdb.Ping().Result()
	log.Println(pong, err)
}
