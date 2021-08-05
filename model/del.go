package model

import "log"
import "github.com/go-redis/redis"

var Redisdb *redis.Client

const RedisKeyNull = redis.Nil //结果为空

func Del() {
	Redisdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := Redisdb.Ping().Result()
	log.Println(pong, err)
}
func Init() {
	if err := Redisdb.SetNX(GetArticleIDKey(), 1, 0).Err(); err != nil {
		panic(err)
	}
}
