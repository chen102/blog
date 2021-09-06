package model

import (
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primary_key"`
	CreatedAt   time.Time `json:"omitempty"`
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index" json:"omitempty"`
	UserID      uint       //用户外键
	UserName    string     `gorm:"-"`
	RevUserName string     `gorm:"-"`
	Content     string
	ArticleID   uint  //文章外键
	FCommentID  int   //父评论ID
	SubComments []int `gorm:"-"`
	Stat        uint  `gorm:"-"`
}

//构建评论树，孩子表示法
func BuildCommentTree(comments []Comment) ([]Comment, []int) {
	var headcomments []int //楼主位置
	visit := make(map[int]int, len(comments))
	for k, comment := range comments {
		visit[int(comment.ID)] = k //记录每条评论记录的下标
		if comment.FCommentID < 0 {
			headcomments = append(headcomments, k)
		}
	}
	for k, comment := range comments {
		if comment.FCommentID >= 0 {
			//若有父亲，取父亲所在地方加入该条记录
			comments[visit[comment.FCommentID]].SubComments = append(comments[visit[comment.FCommentID]].SubComments, k)
		}
	}
	return comments, headcomments
}
