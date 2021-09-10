package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"

	"github.com/gin-gonic/gin"
)

type DeleteService struct {
	DeleteId  uint `form:"DeleteId" json:"DeleteId" binding:"required"`
	ArticleId uint `form:ArticleId json:"ArticleId" binding:"omitempty"` //删除评论时，需要此字段
	Type      bool `form:"Comment" json:"Comment" binding:"omitempty"`
}

func (service *DeleteService) Delete(c *gin.Context) serializer.Response {
	me := model.GetcurrentUser(c)
	if !service.Type { //删除文章
		if !db.ExistArticle(service.DeleteId) {
			return serializer.BuildResponse("没有此文章")
		}
		if err, errorcode := deleteArticle(me.ID, service.DeleteId); err != nil {
			return serializer.Err(errorcode, err)
		}
	} else { //删除评论
		if !db.ExistComment(service.DeleteId) || !db.ExistArticle(service.ArticleId) {
			return serializer.BuildResponse("没有此评论")
		}

		if err, errcode := deleteComment(me.ID, service.DeleteId, service.ArticleId); err != nil {
			return serializer.Err(errcode, err)
		}
	}
	return serializer.BuildResponse("删除成功")
}
func deleteComment(userid, commentid, articleid uint) (err error, errcode int) {
	//更DB
	if err := db.DeleteComment(commentid); err != nil {
		return err, serializer.MysqlErr
	}
	//删缓存
	if err := redis.DeleteComment(articleid); err != nil {
		return err, serializer.RedisErr
	}
	return nil, 0
}
func deleteArticle(userid, articleid uint) (err error, errcode int) {
	//更DB
	if err := db.DeleteArticle(articleid); err != nil {
		return err, serializer.MysqlErr
	}
	//删缓存
	if err := redis.DeleteArticle(userid, articleid); err != nil {
		return err, serializer.RedisErr
	}
	return nil, 0
}
