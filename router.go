package main

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
		v1.GET("list", api.ArticleList)
	}
	return r
}
