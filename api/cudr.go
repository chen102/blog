package api

import (
	"blog/serializer"
	"blog/service"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}
func ArticleList(c *gin.Context) {
	var service service.ArticleSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, ErrorResponse(err))
	} else {
		res := service.ArticleList()
		c.JSON(200, res)

	}
}
func AddArticle(c *gin.Context) {
	var service service.ArticleAddSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, ErrorResponse(err))
	} else {
		res := service.AddArticle()
		c.JSON(200, res)

	}

}
func DeleteArticle(c *gin.Context) {
	var service service.ArticleSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, ErrorResponse(err))
	} else {
		res := service.DeleteArticle()
		c.JSON(200, res)

	}

}
func ShowArticle(c *gin.Context) {
	var service service.ArticleSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, ErrorResponse(err))
	} else {
		res := service.ShowArticle()
		c.JSON(200, res)

	}

}
func UpdateArticle(c *gin.Context) {
	var service service.ArticleSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, ErrorResponse(err))
	} else {
		res := service.UpdateArticle()
		c.JSON(200, res)

	}

}
func ErrorResponse(err error) serializer.Response {
	return serializer.Response{Error: err.Error()}
}
