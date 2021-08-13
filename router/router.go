package router

import (
	"blog/api"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("ping", api.Ping)
	v1 := r.Group("article/manage")
	{
		v1.POST("add", api.AddArticle)
		v1.DELETE("delete", api.DeleteArticle)
		v1.POST("update", api.UpdateArticle)
		v1.POST("show", api.ShowArticle) //传ID了不能用GET
		v1.POST("list", api.ArticleList)
		v1.POST("comment", api.CommentArticle)
		v1.POST("comment/list", api.ShowArticleComment)
		v1.POST("statcomment", api.StatComment)
		//v1.POST("statarticle", api.StatArticle)
		//v1.POST("subcomment")
		//v1.POST("rank")
		//v1.POST("getcomment")
		//v1.POST("getsubcomment")
		//v1.POST("mytags")
		//v1.POST("tagarticle")
		//v1.POST("subscribe") //发布订阅
	}
	return r
}
