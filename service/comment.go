package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"github.com/gin-gonic/gin"
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
	data, err := redis.ShowCommentListCache(service.ArticleId, service.Offset, service.Count)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil {
		comments, err := db.CommentList(service.ArticleId)
		if err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		if err := redis.WriteCommentListCache(service.ArticleId, comments); err != nil {

			return serializer.Err(serializer.RedisErr, err)
		}
		commentTree, roots := model.BuildCommentTree(comments)
		if len(roots) > int(service.Count) {

			return serializer.BuildCommentListResponse(commentTree, roots[:service.Count])
		}
		return serializer.BuildCommentListResponse(commentTree, roots)

	}
	commentTree, roots := model.BuildCommentTree(data)
	return serializer.BuildCommentListResponse(commentTree, roots)
}
