package model

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"omitempty"`
	UpdatedAt time.Time  `json:"omitempty"`
	DeletedAt *time.Time `sql:"index" json:"omitempty"`
	UserName  string
	Account   string `json:"omitempty"`
	Password  string `json:"omitempty"`
	//一对多全部设置外键约束
	Articles    []Article `json:"omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //一个用户有0或多个文章
	Stat        []Article `json:"omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //一个用户对0或多个文章点赞
	Followers   []User    `json:"omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //一个用户对0或多个用户关注
	Fans        []User    `json:"omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FollowerNum uint      `gorm:"-"`         //关注数
	FansNum     uint      `gorm:"-"`         //粉丝数
	Briefly     string    `json:"omitempty"` //简介

	//默认的外键名是拥有者的类型名加上其主键字段名即UserID，也可以重写外键gorm:"foreignKey:xxx
}
type XUser struct {
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

//获取当前登录用户
func GetcurrentUser(c *gin.Context) *User {
	user, _ := c.Get("user") //获取当前用户
	if u, ok := user.(*User); ok {
		return u
	}
	return nil
}
