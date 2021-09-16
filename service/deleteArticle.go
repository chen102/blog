package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"github.com/gin-gonic/gin"
	"log"
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

//评论表存父节点原因
//因为如果存子节点的话，每次添加评论，都需要修改父评论的记录，导致添加评论效率低下，只存父节点的代价是查询和删除都要链式的查找子评论
func deleteComment(userid, commentid, articleid uint) (err error, errcode int) {
	//获取删除的评论及子评论列表
	commentids, err := db.SubCommentid(commentid)
	if err != nil {
		return err, serializer.MysqlErr
	}
	for _, v := range commentids {
		log.Println(v)
	}
	//更DB
	if err := db.DeleteComment(commentids); err != nil {
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
	if err := redis.DeleteArticle(userid); err != nil {
		return err, serializer.RedisErr
	}
	return nil, 0
}
