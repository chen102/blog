package redis

import (
	"blog/model"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func ShowDynamicCache(userid, offset, count uint) ([]model.Article, error) {
	ok, err := model.RedisReadDB.Exists(UserDynamicKey(userid)).Result()
	if err != nil && err != model.RedisNil {
		return nil, err
	} else if ok == 0 {
		return nil, model.RedisNil
	}
	ids, err := model.RedisReadDB.SMembers(UserDynamicKey(userid)).Result()
	if err != nil {
		return nil, err
	}
	get := GetSort(ArticleStringIdKey("*"), true, "UserName", "Title", "UpdatedAt", "Stat", "Tags") //排序后显示的字段
	sortargs := SortArgs("", int64(offset), int64(count), get, "DESC", false)
	data, err := model.RedisWriteDB.Sort(UserDynamicKey(userid), sortargs).Result()
	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, errors.New("没拿到排序结果")
	}
	article, err := model.ArticleList(data)
	if err != nil {
		return nil, err
	}
	//刷新缓存时间
	for _, id := range ids {
		if err := model.RedisWriteDB.Expire(ArticleStringIdKey(id), 1*time.Hour).Err(); err != nil {
			return nil, err
		}
	}
	return article, nil
}
func WriteDynamicCache(userid uint, articles []model.Article) error {
	transactional := func(tx *redis.Tx) error {
		for _, article := range articles {
			//if err := tx.SAdd(UserDynamicKey(userid), article.ID).Err(); err != nil {
			//return err
			//}
			//exist, err := tx.Exists(ArticleIdKey(article.ID)).Result()
			//if err != nil {
			//return err
			//} else if exist == 1 {
			//tx.Expire(ArticleIdKey(article.ID), 1*time.Hour)
			//continue
			//}
			if err := tx.HMSet(ArticleIdKey(article.ID), model.StructToMap(article)).Err(); err != nil {
				return err
			}
			tx.Expire(ArticleIdKey(article.ID), 1*time.Hour)
			//tx.Expire(UserDynamicKey(userid), 10*time.Second)
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, UserDynamicKey(userid)); err != nil { //保证并发安全
		return err
	}

	return nil
}
