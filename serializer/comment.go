package serializer

import (
	"blog/model"
	"log"
	"strconv"
	"time"
)

type Comment struct {
	CommentID   uint
	UserID      uint
	UserName    string
	RevUserName string
	Content     string
	Time        time.Time
	Stat        uint
	SubComment  bool
}

var resComment [][]Comment

func BuildCommentList(comments []model.Comment, roots []int) [][]Comment {
	resC := make([][]Comment, len(roots))
	for i, _ := range resC {
		resC[i] = make([]Comment, 0)
	}
	resComment = resC
	for i, index := range roots {
		comment := Comment{

			CommentID: comments[index].ID,
			UserID:    comments[index].UserID,
			UserName:  comments[index].UserName,
			Content:   comments[index].Content,
			Time:      comments[index].UpdatedAt,
			Stat:      comments[index].Stat,
		}
		resComment[i] = append(resComment[i], comment)
		RangeSubComment(comments, index, i) //index标记楼主在comments的位置，i标记楼主在resComment的位置
	}
	for _, v := range resComment {
		for _, x := range v {
			log.Println(x)
		}
	}
	return resComment
}
func BuildCommentListResponse(comments []model.Comment, roots []int) Response {
	return Response{
		Data: BuildCommentList(comments, roots),
		Msg:  strconv.Itoa(len(resComment)) + " comment Display Succ!",
	}
}
func RangeSubComment(comments []model.Comment, index, i int) {
	for _, v := range comments[index].SubComments {
		comment := Comment{

			CommentID:   comments[v].ID,
			UserID:      comments[v].UserID,
			UserName:    comments[v].UserName,
			Content:     comments[v].Content,
			RevUserName: comments[index].UserName,
			Time:        comments[v].UpdatedAt,
			Stat:        comments[v].Stat,
			SubComment:  true, //表示这是子评论
		}
		resComment[i] = append(resComment[i], comment)
		RangeSubComment(comments, v, i)

	}
}
