package redis

import (
	"blog/model"
	//"errors"
	"github.com/go-redis/redis"
)

func WriteStatCache(stat model.Stat, cancestat bool) error {
	transactional := func(tx *redis.Tx) error {
		//这里文章一定是在redis中的，因为点赞功能需要先点开文章，除非缓存文章的时候出了问题
		//点赞
		if !cancestat {

			if err := tx.HIncrBy(ArticleIdKey(stat.ArticleID), "Stat", 1).Err(); err != nil {
				return err
			}

			if err := tx.LPush(UserStatQueueKey(), UserStatQueueValue(stat.UserID, stat.ArticleID)).Err(); err != nil {
				return err
			}
		} else { //取消点赞
			if err := tx.HIncrBy(ArticleIdKey(stat.ArticleID), "Stat", -1).Err(); err != nil {
				return err
			}

			if err := tx.LPush(UserCancesStatQueueKey(), UserStatQueueValue(stat.UserID, stat.ArticleID)).Err(); err != nil {
				return err
			}

		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticleIdKey(stat.ArticleID)); err != nil { //保证并发安全
		return err
	}
	return nil
}
