package model

import ()

//用户点赞那些文章，可以计算文章点赞量，异步的存入文章表
type StatType uint

const (
	StatArticle StatType = iota
	StatComment
	StatSubComment
)

type Stat struct {
	ID     uint `gorm:"primary_key"`
	Type   StatType
	UserID uint
	StatID uint
	State  bool //用于取消点赞
}
