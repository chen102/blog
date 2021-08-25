package router

import (
	"blog/api"
	"blog/session"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("ping", api.Ping)
	//全局中间件
	r.Use(session.Session("secret"))
	r.Use(session.Cors())
	r.Use(session.CurrentUser())
	v0 := r.Group("/api/v0")
	{
		v0.POST("user/register", api.UserRegister)
		v0.POST("user/login", api.UserLogin)
		v0.POST("article/show", api.ShowArticle) //传ID了不能用GET
		v0.POST("article/list", api.ArticleList)
		v0.POST("article/comment/list", api.ShowArticleComment)
		auth := v0.Group("/")
		auth.Use(session.AuthRequired()) //需要登录的操作
		{
			//auth.DELETE("user/logout", api.UserLogout)
			auth.POST("user/rename", api.UserRename)
			auth.POST("article/add", api.AddArticle)
			auth.DELETE("article/delete", api.DeleteArticle)
			auth.POST("article/update", api.UpdateArticle)
			auth.POST("article/comment", api.CommentArticle)
			auth.POST("article/stat", api.Stat)
		}
	}
	return r
}
