package redis

import (
	"blog/model"
	"time"

	"github.com/go-redis/redis"
)

func ShowUserNameCache(userid uint) (string, error) {
	username, err := model.RedisReadDB.Get(UserIdKey(userid)).Result()
	if err != nil && err != model.RedisNil {
		return "", err
	} else if err == model.RedisNil || username == "" {
		return "", model.RedisNil
	}
	return username, nil

}
func WriteUserNameCache(userid uint, username string) error {
	transactional := func(tx *redis.Tx) error {
		if err := tx.Set(UserIdKey(userid), username, 24*time.Hour).Err(); err != nil {
			return err
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticleStringIdKey(UserIdKey(userid))); err != nil { //保证并发安全
		return err
	}
	return nil
}
func UserIncrementCache(userids []int64) ([]int64, error) {
	return userids, nil
}
