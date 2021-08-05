package service

import (
	"blog/model"
	"blog/serializer"
	"blog/tool"
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
	id, err := model.Redisdb.Get(model.GetArticleIDKey()).Result()
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
	pipe.Incr(model.GetArticleIDKey())
	articleKEY := tool.StrSplicing(model.ArticleIdKey(id))
	if err := pipe.HMSet(articleKEY, article).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	id_int, err := strconv.Atoi(id) //优化:底层使用整数集合
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := pipe.SAdd(model.AuthorArticlesKey(strconv.Itoa(int(service.AuthorId))), id_int).Err(); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	//若设置了标签，使用集合存文章tag，和tag对应的文章
	//可以根据标签找到相应文章
	if service.Tags != nil {
		if err := pipe.SAdd(model.ArticleTagsKey(id), service.Tags).Err(); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		for _, tag := range service.Tags {
			if tag != "" {
				if err := pipe.SAdd(model.TagKey(tag), id_int).Err(); err != nil {
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
	//这里应该先在用户集合里找，有这个用户再去该用户的文章列表找,目前还没做用户模块
	ok, err := model.Redisdb.SIsMember(model.AuthorArticlesKey(strconv.Itoa(int(service.AuthorId))), strconv.Itoa(int(service.ArticleId))).Result() //在该作者所有文章中找到相应id
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	} else if !ok {
		return serializer.BuildResponse("Article ID:" + strconv.Itoa(int(service.ArticleId)) + " NOT EXIST")
	}

	var article model.Article
	data, err := model.Redisdb.HGetAll(model.ArticleIdKey(strconv.Itoa(int(service.ArticleId)))).Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := mapstructure.Decode(data, &article); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildArticleResponse(article, service.ArticleId)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
