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
	AuthorId uint     `form:"AuthorId" json:"AuthorId" binding:"required"`
	Title    string   `form:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content  string   `form:"ArticleContent" json:"ArticleContent" binding:"required"`
	Tags     []string `form:"Tags" json:"Tags" binding:"omitempty"`
}
type ArticleSservice struct {
	AuthorId  uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	ArticleId uint `form:"ArticleId" json:"ArticleId" binding:"required"`
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
		"authorid": service.AuthorId,
		"title":    service.Title,
		"content":  service.Content,
		"time":     tool.ShortTime(),
		"tags":     tool.SliceToString(service.Tags),
	}
	//redis事务
	pipe := model.Redisdb.TxPipeline()
	pipe.Incr("id")
	//使用哈希存文章
	articleKEY := tool.StrSplicing("article:", id) //article:articleID
	if err := pipe.HMSet(articleKEY, article).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	//使用集合存该用户文章id
	articlesKEY := tool.StrSplicing("author:", strconv.Itoa(int(service.AuthorId)), ":", "articles") //author:authorid:articles 用户所有文章
	id_int, err := strconv.Atoi(id)                                                                  //优化:底层使用整数集合
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := pipe.SAdd(articlesKEY, id_int).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	//若设置了标签，使用集合存文章tag，和tag对应的文章
	//可以根据标签找到相应文章
	//找属于多标签的文章，只需将tag:xx:article取交集即可
	if service.Tags != nil {
		tagKEY := tool.StrSplicing(articleKEY, ":", "tags") //article:articleID:tags
		if err := pipe.SAdd(tagKEY, service.Tags).Err(); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		for _, tag := range service.Tags {
			if tag != "" {
				tagArticle := tool.StrSplicing("tag:", tag, ":", "article") //tag:xx:article
				if err := pipe.SAdd(tagArticle, id_int).Err(); err != nil {
					return serializer.Err(serializer.RedisErr, err)
				}

			}
		}
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

	authorArticles := tool.StrSplicing("author:", strconv.Itoa(int(service.AuthorId)), ":", "articles")
	ok, err := model.Redisdb.SIsMember(authorArticles, strconv.Itoa(int(service.ArticleId))).Result() //在该作者所有文章中找到相应id
	if err != nil {
		fmt.Println("DEBUG1")
		return serializer.Err(serializer.RedisErr, err)
	} else if !ok {
		return serializer.BuildResponse("Article ID:" + strconv.Itoa(int(service.ArticleId)) + " NOT EXIST")
	}

	var article model.Article
	redisKEY := tool.StrSplicing("article:", strconv.Itoa(int(service.ArticleId)))
	data, err := model.Redisdb.HGetAll(redisKEY).Result()
	if err != nil {
		fmt.Println("DEBUG2")
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := mapstructure.Decode(data, &article); err != nil {
		fmt.Println("DEBUG3")
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildArticleResponse(article, service.ArticleId)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
