package model

import (
	"time"
)

type Article struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"omitempty"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"omitempty"`
	Title     string
	UserID    uint   //用户的外键
	UserName  string `gorm:"-" map:"omitempty" `
	Stat      uint
	Tags      string `mapstructure:",omitempty"`
	Content   string
}
