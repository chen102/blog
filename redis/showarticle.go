package redis

import (
	"blog/model"
	//"errors"
	//"fmt"
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
	if exist == 0 {
		return nil, model.RedisNil
	}
	data, err := model.RedisReadDB.HGetAll(ArticleIdKey(artid)).Result()
	if err != nil {
		return nil, err
	}
	model.RedisWriteDB.Expire(ArticleIdKey(artid), 1*time.Hour) //刷新存活
	return data, nil

	//if err != nil && err != RedisNil {
	//return nil, err
	//}
	//return data, nil
	//}
	//func ListArticle(uid, offset, count uint, rank bool) ([]string, error) {

	//if count == 0 {
	//count = 5
	//}
	//articlenum, err := Redisdb.SCard(AuthorArticlesKey(uid)).Result()
	//if err != nil && err != RedisNil {
	//return nil, err
	//}
	//if articlenum <= int64(offset) { 当偏移量大于总文章数时，后面返回空
	//return nil, errors.New("is null")
	//}
	//get := GetSort(ArticleIdStringKey(strconv.Itoa(int(uid)), "*"), true, "title", "time", "stat", "tags")
	//var sortargs *redis.Sort
	//if rank { 按点赞排序

	//sortargs = SortArgs(tool.StrSplicing(ArticleIdStringKey(strconv.Itoa(int(uid)), "*"), "->stat"), int64(offset), int64(count), get, "DESC", false)
	//} else { 默认按时间排序
	//sortargs = SortArgs("", int64(offset), int64(count), get, "DESC", false)
	//这里其实应用用也by，拿文章的发布时间来进行排序，但是文章id是系统分配的，id大的一定是后发布的,所以这里直接用key来排序了,后序可以加个按时间排序

	//}

	//res, err := Redisdb.Sort(AuthorArticlesKey(uid), sortargs).Result()
	//if err != nil {
	//return nil, err
	//}
	return nil, nil
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
