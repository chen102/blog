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
	var service service.ArticleListservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.ArticleList(c)
		c.JSON(200, res)

	}
}

func AddArticle(c *gin.Context) {
	var service service.ArticleAddSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.AddArticle(c)
		c.JSON(200, res)

	}

}
func DeleteArticle(c *gin.Context) {
	var service service.ArticleService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.DeleteArticle()
		c.JSON(200, res)

	}

}
func ShowArticle(c *gin.Context) {
	var service service.ArticleService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.ShowArticle()
		c.JSON(200, res)

	}

}
func UpdateArticle(c *gin.Context) {
	var service service.ArticleService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.UpdateArticle()
		c.JSON(200, res)

	}

}
func Stat(c *gin.Context) {
	var service service.StatService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Stat(c)
		c.JSON(200, res)
	}
}
func Comment(c *gin.Context) {
	var service service.CommentService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Comment(c)
		c.JSON(200, res)
	}

}
func ShowCommentList(c *gin.Context) {
	var service service.ArticleCommentListservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.CommentList()
		c.JSON(200, res)
	}

}
