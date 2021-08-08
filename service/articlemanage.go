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
type ArticleListservice struct {
	AuthorId uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	Offset   uint `form:"ArticleOffset" json:"ArticleOffset" binding:"omitempty"`
	Count    uint `form:"ArticleCount" json:"ArticleCount" binding:"omitempty"`
}

func (service *ArticleListservice) ArticleList() serializer.Response {
	//这里应该也是一样，先在用户集合了找
	if service.Count == 0 {
		service.Count = 5
	}
	articlenum, err := model.Redisdb.SCard(model.AuthorArticlesKey(strconv.Itoa(int(service.AuthorId)))).Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	fmt.Printf("%T", articlenum)
	if articlenum <= int64(service.Offset) { //当偏移量大于总文章数时，后面返回空
		return serializer.BuildResponse("is Null")
	}
	get := model.GetSort(model.ArticleIdKey("*"), true, "title", "time", "tags")
	sortargs := model.SortArgs("", int64(service.Offset), int64(service.Count), get, "DESC", false)
	//这里其实应用用by，拿文章的发布时间来进行排序，但是文章id是系统分配的，id大的一定是后发布的,所以这里直接用key来排序了

	res, err := model.Redisdb.Sort(model.AuthorArticlesKey(strconv.Itoa(int(service.AuthorId))), sortargs).Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	//Sort返回的结果为string，将string转为多个文章模型进行响应
	article := make([]model.Article, len(res)/4) //4个string为一个article,分别是id,title,time,tags
	id := 0
	for i := 0; i < len(res); i++ {
		if i != 0 && i%4 == 0 {
			id++
		}
		switch i % 4 {
		case 0:
			articleid, err := strconv.Atoi(res[i])
			if err != nil {
				return serializer.Err(serializer.StrconvErr, err)
			}
			article[id].ArticleId = uint(articleid)
		case 1:
			article[id].Title = res[i]
		case 2:
			article[id].Time = res[i]
		case 3:
			if res[i] != "" {

				article[id].Tags = res[i]
			}
		}
	}
	return serializer.BuildArticleListResponse(article)

}
func (service *ArticleAddSservice) AddArticle() serializer.Response {
	id, err := model.Redisdb.Get(model.GetArticleIDKey()).Result()
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}

	article := map[string]interface{}{
		"title":   service.Title,
		"content": service.Content,
		"time":    tool.ShortTime(),
	}
	if service.Tags != nil {
		article["tags"] = tool.SliceToString(service.Tags)
	}
	//redis事务
	pipe := model.Redisdb.TxPipeline()
	pipe.Incr(model.GetArticleIDKey())
	if err := pipe.HMSet(model.ArticleIdKey(id), article).Err(); err != nil {
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
	article.ArticleId = service.ArticleId
	article.AuthorId = service.AuthorId
	return serializer.BuildArticleResponse(article)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
