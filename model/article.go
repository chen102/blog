package model

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title    string
	UserID   uint   //用户的外键
	UserName string `gorm:"-"`
	Stat     uint
	Tags     string
	Content  string
}
