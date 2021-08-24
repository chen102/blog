package model

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title    string
	UserID   uint   //用户的外键
	UserName string `gorm:"-" map:"omitempty" json:"omitempty"`
	Stat     uint
	Tags     string `mapstructure:",omitempty"`
	Content  string
}
