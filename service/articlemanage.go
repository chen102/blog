package service

import (
	"blog/model"
	"blog/serializer"
	"blog/tool"
	"fmt"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type ArticleAddSservice struct {
	Title   string `from:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content string `from:"ArticleContent" json:"ArticleContent" binding:"required"`
}
type ArticleSservice struct {
	Id uint `from:"id" json:"id" binding:"required"`
}

func (service *ArticleSservice) ArticleList() serializer.Response {
	//ids, err := model.Redisdb.SMembers("articles").Result()
	//if err != nil {
	//return serializer.Err(err)
	//}
	return serializer.Response{}
}
func (service *ArticleAddSservice) AddArticle() serializer.Response {
	id, err := model.Redisdb.Get("id").Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}

	article := map[string]interface{}{
		"title":   service.Title,
		"content": service.Content,
		"time":    tool.ShortTime(),
	}
	//redis事务
	pipe := model.Redisdb.TxPipeline()
	pipe.Incr("id")
	redisKEY := tool.StrSplicing("article:", id)
	//使用hash存文章
	if err := pipe.HMSet(redisKEY, article).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	//使用set存文章id
	if err := pipe.SAdd("articles", id).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if _, err := pipe.Exec(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("Article ID:" + id + " ADD Succ！")
}
func (service *ArticleSservice) DeleteArticle() serializer.Response {
	return serializer.Response{}

}
func (service *ArticleSservice) ShowArticle() serializer.Response {
	var article model.Article
	redisKEY := tool.StrSplicing("article:", strconv.Itoa(int(service.Id)))
	fmt.Println(redisKEY)
	data, err := model.Redisdb.HGetAll(redisKEY).Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := mapstructure.Decode(data, &article); err != nil {
		fmt.Println("ERROR:", err)
		return serializer.Err(serializer.RedisErr, err)
	}
	fmt.Println(article)
	return serializer.BuildArticleResponse(article, service.Id)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
