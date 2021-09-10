package model

import (
	"blog/tool"
	"strconv"
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
	RootID      int   //楼主
	SubComments []int `gorm:"-" json:"omitempty"`
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
			//方便序列化
			//if comments[visit[comment.FCommentID]].SubCommentsString == "" {
			//comments[visit[comment.FCommentID]].SubCommentsString = strconv.Itoa(k)
			//} else {
			//comments[visit[comment.FCommentID]].SubCommentsString = comments[visit[comment.FCommentID]].SubCommentsString + "," + strconv.Itoa(k)
			//}
			//}
		}
	}
	return comments, headcomments
}
func CommentRank(data []string) ([]Comment, error) {
	comments := make([]Comment, len(data)/8)
	id := 0
	for i := 0; i < len(data); i++ {
		if i != 0 && i%8 == 0 {
			id++
		}
		switch i % 8 {
		case 0:
			commentid, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}
			comments[id].ID = uint(commentid)
		case 1:
			userid, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}

			comments[id].UserID = uint(userid)
		case 2:
			comments[id].UserName = data[i]
		case 3:
			comments[id].RevUserName = data[i]
		case 4:
			comments[id].UpdatedAt = tool.StringToTime(data[i])
		case 5:
			stat, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}
			comments[id].Stat = uint(stat)
		case 6:
			comments[id].Content = data[i]
		case 7:
			fcommentid, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}

			comments[id].FCommentID = fcommentid
		}
	}
	return comments, nil
}
