package router

import (
	"blog/api"
	_ "blog/docs"
	"blog/session"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("ping", api.Ping)
	//全局中间件
	r.Use(session.Session("secret"))
	r.Use(session.Cors())
	r.Use(session.CurrentUser())
	//swagger := &ginSwagger.Config{
	//URL: "http://localhost:3000/swagger/doc.json",
	//}
	v0 := r.Group("/api/v0")
	{
		v0.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v0.POST("user/register", api.UserRegister)
		v0.POST("user/login", api.UserLogin)
		v0.GET("article/show", api.ShowArticle)
		auth := v0.Group("/")
		auth.Use(session.AuthRequired()) //需要登录的操作
		{
			//auth.DELETE("user/logout", api.UserLogout)
			auth.POST("article/comment/list", api.ShowCommentList)
			auth.POST("article/delete", api.Delete)
			auth.POST("article/comment", api.Comment)
			auth.POST("follower/dynamic", api.UserDynamicList)
			auth.POST("follower", api.UserFollowerUser)
			auth.POST("follower/list", api.UserFollowerList)
			auth.POST("user/like", api.UserArticlesLike)
			auth.POST("article/list", api.ArticleList)
			auth.POST("article/add", api.AddArticle)
			auth.POST("article/stat", api.Stat)
		}
	}
	return r
}
