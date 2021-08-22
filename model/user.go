package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	UserName string
	Account  string
	Password string
	Articles []Article //一个用户有0或多个文章
	Stat     []Article //一个用户对0或多个文章点赞

	//默认的外键名是拥有者的类型名加上其主键字段名即UserID，也可以重写外键gorm:"foreignKey:xxx
}

func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}
func (user *User) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

//获取当前登录用户的ID
func GetcurrentID(c *gin.Context) *User {
	user, _ := c.Get("user") //获取当前用户
	if u, ok := user.(*User); ok {
		return u
	}
	return nil
}
