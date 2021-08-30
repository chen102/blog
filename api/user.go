package api

import (
	"blog/serializer"
	"blog/service"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var service service.UserRegisterService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Register()
		c.JSON(200, res)

	}
}
func UserLogin(c *gin.Context) {
	var service service.UserLoginService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Login(c)
		c.JSON(200, res)

	}
}
func UserRename(c *gin.Context) {
	var service service.UpdateUsernameService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.UpdateUsername(c)
		c.JSON(200, res)

	}

}
func UserArticlesLike(c *gin.Context) {
	var service service.UserStatListService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.UserStatArticlesList(c)
		c.JSON(200, res)

	}

}
