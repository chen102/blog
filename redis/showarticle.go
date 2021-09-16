package redis

import (
	"blog/model"
	"blog/tool"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func WriteArticleCache(article map[string]interface{}) error {
	articleid := strconv.FormatFloat(article["ID"].(float64), 'f', 0, 64)
	transactional := func(tx *redis.Tx) error {
		if err := tx.HMSet(ArticleStringIdKey(articleid), article).Err(); err != nil {
			return err
		}
		tx.Expire(ArticleStringIdKey(articleid), 1*time.Hour) //1小时存活
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticleStringIdKey(articleid)); err != nil { //保证并发安全
		return err
	}
	return nil

}
func ShowArticleCache(artid uint) ([]model.Article, error) {
	exist, err := model.RedisReadDB.Exists(ArticleIdKey(artid)).Result()
	if err != nil && err != model.RedisNil {
		return nil, err
	} else if exist == 0 {
		return nil, model.RedisNil
	}
	data, err := model.RedisReadDB.HGetAll(ArticleIdKey(artid)).Result()
	if err != nil {
		return nil, err
	}
	model.RedisWriteDB.Expire(ArticleIdKey(artid), 1*time.Hour) //刷新存活
	article := make([]model.Article, 1)
	stat, err := strconv.Atoi(data["Stat"])
	if err != nil {
		return nil, err
	}
	userid, err := strconv.Atoi(data["UserID"])
	if err != nil {
		return nil, err
	}
	article[0].ID = artid
	article[0].Stat = uint(stat)
	article[0].UserID = uint(userid)
	article[0].Content = data["Content"]
	article[0].UserName = data["UserName"]
	article[0].Title = data["Title"]
	article[0].UpdatedAt = tool.StringToTime(data["UpdatedAt"])
	article[0].Tags = data["Tags"]
	return article, nil
}
func ShowArticleListCache(userid, offset, count uint, rank bool) ([]model.Article, error) {
	var res []string
	if rank {
		//rank功能待做，思路:需使用有序集合，每次点赞和查看文章时，增加文章的热度(一定时间内),然后按热度排序
		//sortargs = SortArgs(tool.StrSplicing(ArticleIdStringKey(strconv.Itoa(int(uid)), "*"), "->stat"), int64(offset), int64(count), get, "DESC", false)
	} else {
		exist, err := model.RedisReadDB.Exists(ArticlesListKey(userid)).Result()
		if err != nil && err != model.RedisNil {
			return nil, err
		} else if exist == 0 {
			return nil, model.RedisNil
		}
		get := GetSort(ArticleStringIdKey("*"), true, "UserName", "Title", "UpdatedAt", "Stat", "Tags") //排序后显示的字段
		sortargs := SortArgs("", int64(offset), int64(count), get, "DESC", false)                       //按时间排序：直接按key排序,也就是ID，ID大的一定后发布
		res, err = model.RedisWriteDB.Sort(ArticlesListKey(userid), sortargs).Result()
		//这里应该把排序的结果保存起来，下次直接查
		if err != nil {
			return nil, err
		}
	}
	if res == nil {
		return nil, errors.New("没拿到排序结果")
	}
	article, err := model.ArticleList(res)
	if err != nil {
		return nil, err
	}
	return article, nil
}
func WriteArticleListCach(userid uint, articles []model.Article) error {
	//写入用户文章合集
	//写入用户文章
	//取出文章的ID
	transactional := func(tx *redis.Tx) error {

		for _, v := range articles {
			if err := tx.HMSet(ArticleIdKey(v.ID), model.StructToMap(v)).Err(); err != nil {
				return err
			}
			tx.Expire(ArticleIdKey(v.ID), 1*time.Hour)
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticlesListKey(userid)); err != nil { //保证并发安全
		return err
	}
	return nil
}
func DeleteArticle(userid uint) error {
	//删除文章缓存 1.删动态索引 2.删文章索引
	if err := model.RedisWriteDB.Del(UserDynamicKey(userid)); err != nil {
		return nil
	}
	if err := model.RedisWriteDB.Del(ArticlesListKey(userid)).Err(); err != nil {
		return nil
	}
	return nil
}
func ArticleIncrementCache(userid uint, articleids []int64) ([]int64, error) {
	transactional := func(tx *redis.Tx) error {
		for k, id := range articleids {
			if err := tx.SAdd(UserDynamicKey(userid), id).Err(); err != nil {
				return err
			}
			exist, err := tx.Exists(ArticleIdKey(uint(id))).Result()
			if err != nil && err != model.RedisNil {
				return err
			} else if exist == 1 {
				if err := tx.Expire(ArticleIdKey(uint(id)), 1*time.Hour).Err(); err != nil {
					return err
				}
				articleids[k] = -1
			}
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticlesListKey(userid)); err != nil { //保证并发安全
		return nil, err
	}
	return articleids, nil
}
