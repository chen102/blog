package redis

import (
	"blog/model"
	"errors"
)

func Stat(stat model.Stat, article map[string]interface{}) error {
	exis, err := model.Redisdb.Exists(UserStatArticleKey(stat.UserID, stat.ArticleID)).Result()
	if err != nil {
		return err
	}
	if exis == 1 {
		return errors.New("别点了，点过了")
	}
	pipe := model.Redisdb.Pipeline()
	pipe.HMSet(UserStatArticleKey(stat.UserID, stat.ArticleID), article)
	//更新用户点赞列表同时防止重复点赞
	//点赞数+1
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}
