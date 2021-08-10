package redis

import (
	. "blog/model"
	"errors"
	"strconv"
)

func ShowArticle(uid, artid uint) (interface{}, error) {

	//这里应该先在用户集合里找，有这个用户再去该用户的文章列表找,目前还没做用户模块
	ok, err := Redisdb.SIsMember(AuthorArticlesKey(uid), artid).Result() // 在该作者所有文章中找到相应id
	if err != nil && err != RedisNil {
		return nil, err
	} else if !ok {
		return nil, errors.New("artid is not exist")
	}
	data, err := Redisdb.HGetAll(ArticleIdKey(uid, artid)).Result()
	if err != nil && err != RedisNil {
		return nil, err
	}
	return data, nil
}
func ListArticle(uid, offset, count uint) ([]string, error) {

	if count == 0 {
		count = 5
	}
	articlenum, err := Redisdb.SCard(AuthorArticlesKey(uid)).Result()
	if err != nil && err != RedisNil {
		return nil, err
	}
	if articlenum <= int64(offset) { //当偏移量大于总文章数时，后面返回空
		return nil, errors.New("is null")
	}
	get := GetSort(ArticleIdStringKey(strconv.Itoa(int(uid)), "*"), true, "title", "time", "tags")
	sortargs := SortArgs("", int64(offset), int64(count), get, "DESC", false)
	//这里其实应用用by，拿文章的发布时间来进行排序，但是文章id是系统分配的，id大的一定是后发布的,所以这里直接用key来排序了,后序可以加个按时间排序

	res, err := Redisdb.Sort(AuthorArticlesKey(uid), sortargs).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}
func ShowComment(uid, artid uint) ([]string, error) {
	commentnum, err := Redisdb.ZCard(ArticleCommentRankKey(uid, artid)).Result()
	if err != nil {
		return nil, err
	}
	if commentnum == 0 {
		return nil, errors.New("There are no comments at this time")
	}
	if commentnum < 5 {
		comment := make([]string, commentnum)
		commentids, err := Redisdb.ZRange(ArticleCommentRankKey(uid, artid), 0, 5).Result()
		if err != nil {
			return nil, err
		}
		for k, commentid := range commentids {
			comment[k], err = Redisdb.HGet(ArticleCommentIDStringKey(strconv.Itoa(int(uid)), strconv.Itoa(int(artid)), commentid), "comment:0").Result()

		}
		return comment, nil

		//} else {

		//}
		////默认显示5条

	}
	return nil, nil
}
