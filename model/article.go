package model

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title   string
	UserID  uint
	User    User `gorm:"ForeignKey:UserID"` //使用UserID作为外键
	Stat    uint
	Tags    string
	Content string
}
