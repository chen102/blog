package service

import (
	"blog/model"
	"blog/model/db"
	"blog/serializer"
	"github.com/gin-gonic/gin"
	"log"
)

type CommentService struct {
	ArticleService
	CommentID uint   `form:"Rev" json:"Rev" binding:"omitempty"`
	Content   string `form:"Content" json:"Content" binding:"required"`
}
type ArticleCommentListservice struct {
	ArticleService

	CommentID uint `form:"Rev" json:"Rev" binding:"omitempty"`
	Paginationservice
}

func (service *CommentService) Comment(c *gin.Context) serializer.Response {
	if !db.ExistArticle(service.ArticleId) {
		return serializer.BuildResponse("没有此文章")
	}
	me := model.GetcurrentUser(c)

	comment := model.Comment{
		UserID:    me.ID,
		ArticleID: service.ArticleId,
		Content:   service.Content,
	}
	if service.CommentID == 0 {
		comment.FCommentID = -1
	} else {

		if !db.ExistComment(service.CommentID) {

			return serializer.BuildResponse("没有此评论")
		}
		comment.FCommentID = int(service.CommentID)
	}
	if err := model.DB.Create(&comment).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, err)
	}

	return serializer.BuildResponse("xx")
}
func (service *ArticleCommentListservice) CommentList() serializer.Response {
	if service.Count == 0 {
		service.Count = 10
	}
	comments, err := db.CommentList(service.ArticleId)
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	commentTree, roots := model.BuildCommentTree(comments)
	log.Println(roots)
	for _, v := range commentTree {

		log.Println(v)
	}
	return serializer.BuildCommentListResponse(commentTree, roots)
	//return serializer.BuildResponse("xx")
}
