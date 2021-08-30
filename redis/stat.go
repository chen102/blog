package redis

import (
	"blog/model"
	"blog/tool"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

//go的bool不支持转int只能这样了
//异或关系
//存在，点赞 1
//存在，取消点赞 0
//不存在，点赞 1
//不存在，取消点赞 0
func XOR(exist, cancestat bool) bool {
	e, c := 0, 0
	if exist {
		e = 1
	}
	if cancestat {
		c = 1
	}
	return func(a, b int) bool {
		return a^b == 1
	}(e, c)
}

//redis事务到底要不要，处理每条错误
//按道理来说，是不用，因为是一起返回结果
func WriteStatCache(stat model.Stat, cancestat bool) error {
	//防止重复点赞
	if ok, err := model.RedisReadDB.SIsMember(UserStatList(stat.UserID), stat.ArticleID).Result(); err != nil {
		return err
	} else if XOR(ok, cancestat) {
		return errors.New("您可别点了")
	}
	//进入点赞(取消)事务
	transactional := func(tx *redis.Tx) error {
		//这里文章一定是在redis中的，因为点赞功能需要先点开文章，除非缓存文章的时候出了问题
		//点赞
		if !cancestat {
			//点赞数+1
			if err := tx.HIncrBy(ArticleIdKey(stat.ArticleID), "Stat", 1).Err(); err != nil {
				return err
			}
			//加入用户点赞列表
			if err := tx.SAdd(UserStatList(stat.UserID), stat.ArticleID).Err(); err != nil {
				return err
			}

			if err := tx.LPush(UserStatQueueKey(), UserStatQueueValue(stat.UserID, stat.ArticleID)).Err(); err != nil {
				return err
			}
		} else { //取消点赞
			if err := tx.HIncrBy(ArticleIdKey(stat.ArticleID), "Stat", -1).Err(); err != nil {
				return err
			}
			//这里不用担心并发的问题，因为这是用户自己的点赞列表，所以直接可以修改缓存
			if err := tx.SRem(UserStatList(stat.UserID), stat.ArticleID); err != nil {
				return nil
			}
			//点赞和取消点赞公用一个队列
			if err := tx.LPush(UserStatQueueKey(), UserCancesStatQueueValue(stat.UserID, stat.ArticleID)).Err(); err != nil {
				return err
			}

		}
		return nil
	}
	//这里不仅要监视文章，还要监视用户点赞列表
	if err := model.RedisWriteDB.Watch(transactional, ArticleIdKey(stat.ArticleID), UserStatList(stat.UserID)); err != nil { //保证并发安全
		return err
	}
	return nil
}

//判断用户点赞列表缓存是否过期
func ExistUserStatList(userid uint) error {
	exist, err := model.RedisReadDB.Exists(UserStatList(userid)).Result()
	if err != nil && err != model.RedisNil {
		return err
	} else if exist == 0 {
		return model.RedisNil
	}
	return nil
}
func ShowUserStatListCache(userid, offset, count uint) ([]model.Article, error) {
	ok, err := model.RedisReadDB.Exists(UserStatList(userid)).Result()
	if err != nil && err != model.RedisNil {
		return nil, err
	} else if ok == 0 {
		return nil, model.RedisNil
	}
	ids, err := model.RedisReadDB.SMembers(UserStatList(userid)).Result()
	if err != nil {
		return nil, err
	}

	model.RedisWriteDB.Expire(UserStatList(userid), 1*time.Hour)
	//set没有分页功能，手动实现下
	page := tool.Pagination(ids, int(offset), int(count))
	if page == nil {
		return nil, nil
	}
	articles := make([]model.Article, len(page))
	for k, id := range page {
		data, err := model.RedisReadDB.HMGet(ArticleStringIdKey(id), "Title", "UpdatedAt", "Stat", "Tags").Result()
		if err != nil {
			return nil, err
		}
		if data != nil {
			log.Println(data)
			articleid, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			stat, err := strconv.Atoi(data[2].(string))
			if err != err {
				return nil, err
			}
			//这代码属实让人想吐
			//纠其原因是因为定义的模型、go-redis输入输出、gorm输入输出的数据类型没有统一,所以老是在转换转换转换(还是太年轻了)
			articles[k].ID = uint(articleid)
			articles[k].Title = data[0].(string)
			articles[k].UpdatedAt = tool.StringToTime(data[1].(string))
			articles[k].Stat = uint(stat)
			if data[3] != nil {

				articles[k].Tags = data[3].(string)
			}

			model.RedisWriteDB.Expire(ArticleIdKey(uint(articleid)), 1*time.Hour)
		}
	}
	return articles, nil

}

//和文章列表缓存一样
//一个文章ID合集(索引)，公用文章缓存
func WriteUserStatListCache(userid uint, articles []model.Article) error {

	transactional := func(tx *redis.Tx) error {
		for k, article := range articles {
			if err := tx.SAdd(UserStatList(userid), article.ID).Err(); err != nil {
				return err
			}
			tx.Expire(UserStatList(userid), 1*time.Hour)
			exist, err := tx.Exists(ArticleIdKey(article.ID)).Result()
			if err != nil {
				return err
			} else if exist == 1 {
				tx.Expire(ArticleIdKey((article.ID)), 1*time.Hour)
				continue
			}
			log.Println(articles[k])
			log.Println(model.StructToMap(articles[k]))
			if err := tx.HMSet(ArticleIdKey(article.ID), model.StructToMap(articles[k])).Err(); err != nil {
				return err
			}
			tx.Expire(ArticleIdKey(article.ID), 1*time.Hour) //1小时存活
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, UserStatList(userid)); err != nil {
		return err
	}

	return nil
}
