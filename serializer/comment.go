package serializer

import (
	"blog/model"
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

var resComment []Comment

func BuildCommentList(comments []model.Comment, roots []int) []Comment {
	resC := make([]Comment, 0)
	resComment = resC
	for _, index := range roots {
		comment := Comment{

			CommentID: comments[index].ID,
			UserID:    comments[index].UserID,
			UserName:  comments[index].UserName,
			Content:   comments[index].Content,
			Time:      comments[index].UpdatedAt,
			Stat:      comments[index].Stat,
		}
		resComment = append(resComment, comment)
		RangeSubComment(comments, index)
	}
	return resComment
}
func BuildCommentListResponse(comments []model.Comment, roots []int) Response {
	return Response{
		Data: BuildCommentList(comments, roots),
		Msg:  strconv.Itoa(len(resComment)) + " comment Display Succ!",
	}
}
func RangeSubComment(comments []model.Comment, index int) {
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

		resComment = append(resComment, comment)
		RangeSubComment(comments, v)

	}
}
