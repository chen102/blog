package redis

import (
	"blog/model"
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
func ShowArticleCache(artid uint) (interface{}, error) {
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
	return data, nil
}
func ShowArticleListCache(userid, offset, count uint, rank bool) ([]string, error) {
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
		get := GetSort(ArticleStringIdKey("*"), true, "Title", "UpdatedAt", "Stat", "Tags") //排序后显示的字段
		sortargs := SortArgs("", int64(offset), int64(count), get, "DESC", false)           //按时间排序：直接按key排序,也就是ID，ID大的一定后发布
		res, err = model.RedisWriteDB.Sort(ArticlesListKey(userid), sortargs).Result()
		//这里应该把排序的结果保存起来，下次直接查
		if err != nil {
			return nil, err
		}
	}
	if res == nil {
		return nil, errors.New("没拿到排序结果")
	}
	return res, nil
}
func WriteArticleListCach(userid uint, articles []model.Article) error {
	//写入用户文章合集
	//写入用户文章
	//取出文章的ID
	transactional := func(tx *redis.Tx) error {

		for _, v := range articles {
			if err := tx.SAdd(ArticlesListKey(userid), v.ID).Err(); err != nil {
				return err
			}
			exist, err := tx.Exists(ArticlesListKey(v.ID)).Result()
			if err != nil {
				return err
			} else if exist == 1 {
				tx.Expire(ArticleIdKey(v.ID), 1*time.Hour) //这里就把之前缓存的但这次没缓存的文章的过期时间强制与这次缓存的过期时间一致)
				continue                                   //若缓存存在，跳过，避免struct转map的开销
			}
			if err := tx.HMSet(ArticleIdKey(v.ID), model.StructToMap(v)).Err(); err != nil {
				return err
			}
			tx.Expire(ArticleIdKey(v.ID), 1*time.Hour)
			tx.Expire(ArticlesListKey(userid), 1*time.Hour)
		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticlesListKey(userid)); err != nil { //保证并发安全
		return err
	}
	return nil
}
func ShowComment(uid, artid uint) ([]string, error) {
	//commentnum, err := Redisdb.ZCard(ArticleCommentRankKey(uid, artid)).Result()
	//if err != nil {
	//return nil, err
	//}
	//if commentnum < 5 {
	//comment := make([]string, commentnum)
	//commentids, err := Redisdb.ZRevRange(ArticleCommentRankKey(uid, artid), 0, commentnum).Result()
	//if err != nil {
	//return nil, err
	//}
	//for k, commentid := range commentids {
	//comment[k], err = Redisdb.HGet(ArticleCommentIDStringKey(strconv.Itoa(int(uid)), strconv.Itoa(int(artid)), commentid), "comment:0").Result()

	//}
	//comment = append(comment, strconv.Itoa(int(commentnum))) //最后一个string判断是否还有评论
	//return comment, nil

	//} else {

	//comment := make([]string, 5)
	//commentids, err := Redisdb.ZRange(ArticleCommentRankKey(uid, artid), 0, 5).Result()
	//if err != nil {
	//return nil, err
	//}
	//for k, commentid := range commentids {
	//comment[k], err = Redisdb.HGet(ArticleCommentIDStringKey(strconv.Itoa(int(uid)), strconv.Itoa(int(artid)), commentid), "comment:0").Result()

	//}

	//comment = append(comment, strconv.Itoa(int(commentnum))) //默认只显示5条评论，并把总评论数返回
	//return comment, nil
	//}
	return nil, nil

}
func ShowAllComment(uid, artid, offset, count uint) ([]string, error) {
	//comment := make([]string, count)
	//stop := count + offset - 1
	//commentids, err := Redisdb.ZRevRange(ArticleCommentRankKey(uid, artid), int64(offset), int64(stop)).Result()
	//if err != nil {
	//return nil, err
	//}
	//for k, commentid := range commentids {
	//comment[k], err = Redisdb.HGet(ArticleCommentIDStringKey(strconv.Itoa(int(uid)), strconv.Itoa(int(artid)), commentid), "comment:0").Result()

	//}
	return nil, nil

}
func GetStat(uid, artid, commentid uint) (uint, error) {
	//stat, err := Redisdb.Get(ArticleCommentStatKey(uid, artid, commentid)).Int()
	//if err != nil && err != RedisNil {
	//return 0, err
	//}
	//return uint(stat), nil
	return 0, nil
}
