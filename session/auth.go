package session

import (
	"blog/model"
	"blog/model/db"
	"blog/serializer"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("userID") //从会话获得用户id
		if uid != nil {
			user, err := db.GetUser(uid)
			if err == nil {
				c.Set("user", &user)
			}
		}
		c.Next() //等待执行其他中间件
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			if _, ok := user.(*model.User); ok {
				c.Next()
				return
			}
		}

		c.JSON(200, serializer.Response{
			Code: 401,
			Msg:  "需要登录",
		})
		c.Abort() //Abort 函数在被调用的函数中阻止后续中间件的执行
	}
}
