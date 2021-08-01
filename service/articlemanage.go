package service

import (
	"blog/model"
	"blog/serializer"
)

type ArticleSservice struct {
	Title   string `from:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content string `from:"ArticleContent" json:"ArticleContent" binding:"required"`
}

func ArticleList() serializer.Response {
	return serializer.Response{}
}
func AddArticle(service *ArticleSservice) serializer.Response {
	//var article model.Article
	//id := model.Redisdb.Get("uid").Result()
	////事务
	//pipe := model.Redisdb.TxPipeline()
	//pipe.Incr("uid")
	//if err := pipe.HSet(); err != nil {
	//return serializer.Err(err)
	//}
	//if _, err := pipe.Exec(); err != nil {
	//return serializer.Err(err)
	//}
	return serializer.Response{}
}
func DeleteArticle() serializer.Response {
	return serializer.Response{}

}
func ShowArticle() serializer.Response {
	return serializer.Response{}

}
func UpdateArticle() serializer.Response {
	return serializer.Response{}
}
