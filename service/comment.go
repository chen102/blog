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
	CommentID  uint   `form:"Rev" json:"Rev" binding:"omitempty"`             //回复评论ID
	LandlordID uint   `form:"LandlordId json:"LandlordId binding:"omitempty"` //楼主ID
	Content    string `form:"Content" json:"Content" binding:"required"`
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
	if service.CommentID == 0 || service.LandlordID == 0 {
		comment.FCommentID = -1
	} else {

		if !db.ExistComment(service.CommentID) || !db.ExistComment(service.LandlordID) {

			return serializer.BuildResponse("没有此评论")
		}
		comment.FCommentID = int(service.CommentID)
		comment.RootID = int(service.LandlordID)
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
		//拿评论ID
		commentids, err := db.CommentID(service.ArticleId)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		//筛选ID
		commentids, err = redis.CommentIncrementCache(service.ArticleId, commentids)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)

		}
		//拿评论信息
		comments, err := db.CommentList(commentids)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteCommentListCache(service.ArticleId, comments); err != nil {

			return serializer.Err(serializer.RedisErr, err)
		}
		//commentTree, roots := model.BuildCommentTree(comments)
		//if len(roots) > int(service.Count) {

		//return serializer.BuildCommentListResponse(commentTree, roots[:service.Count])
		//}
		//return serializer.BuildCommentListResponse(commentTree, roots)
		return serializer.BuildResponse("请再次执行")

	}
	commentTree, roots := model.BuildCommentTree(data)
	return serializer.BuildCommentListResponse(commentTree, roots)
}
