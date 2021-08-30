package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"

	//"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	//"github.com/mitchellh/mapstructure"
	//"encoding/json"
)

type StatService struct {
	//文章、评论、子评论的ID
	StatId uint `form:"StatId" json:"StatId" binding:"required"`
	//0 1 2 对文章、评论、子评论点赞
	StatType model.StatType `form:"StatType" json:"StatType" binding:"omitempty,oneof=1 2"`
	//默认为点赞，加入这个参数为取消点赞
	CancelStat bool `form:"CancelStat" json:"CancelStat" binding:"omitempty"`
}

func (service *StatService) Stat(c *gin.Context) serializer.Response {
	session := sessions.Default(c)
	userid := session.Get("userID")
	//点赞模型
	stat := model.Stat{
		Type:      model.StatArticle,
		UserID:    userid.(uint),
		ArticleID: service.StatId,
	}
	switch service.StatType {
	case 0:
		stat.Type = model.StatArticle
		if err, errcode := StatArticle(stat, service.CancelStat); err != nil {
			return serializer.Err(errcode, err)
		}
	case 1:
		stat.Type = model.StatComment
		if err, errcode := StatComment(stat, service.CancelStat); err != nil {
			return serializer.Err(errcode, err)
		}

	case 2:
		stat.Type = model.StatSubComment
		if err, errcode := StatSubComment(stat, service.CancelStat); err != nil {
			return serializer.Err(errcode, err)
		}
	}
	return serializer.BuildResponse("点赞成功!")
}

//errorCode 返回错误代码，用于上层捕获err类型
//点赞功能设计
//用户点赞，写入redis用户点赞合集，redis定时写入mysql点赞表，mysql计算点赞数，用户点赞列表 写入redis，

func StatArticle(stat model.Stat, cancestat bool) (err error, errcode int) {
	//刷新用户点赞列表
	//这里刷新用户点赞列表目的：用于及时的更新用户点赞列表，因为用户点赞列表是从点赞表查的，而点赞写缓存了，会造成用户点赞后，无法第一时间更新点赞列表(mysql、redis缓存数据一致性)
	//第二个作用是防止用户重复点赞
	if err := redis.ExistUserStatList(stat.UserID); err != nil && err != model.RedisNil {
		return err, serializer.RedisErr
	} else if err == model.RedisNil {
		var ids []model.Stat
		//从点赞表获取该用户点赞文章写入cache
		if err := model.DB.Select("article_id").Where("user_id=? AND stat=?", stat.UserID, 0).Find(&ids).Error; err != nil {
			return err, serializer.MysqlErr
		}
		var articles []model.Article
		for k, id := range ids {
			//计算这篇文章的点赞数
			//if err:=model.DB.Where()
			if err := model.DB.Where("article_id=?", id.ArticleID).First(&articles[k]).Error; err != nil {

				return err, serializer.MysqlErr
			}

		}
		if err := redis.WriteUserStatListCache(stat.UserID, articles); err != nil {

			return err, serializer.MysqlErr
		}

	}

	//点赞写缓存
	if err := redis.WriteStatCache(stat, cancestat); err != nil {
		return err, serializer.RedisErr
	}
	return nil, 0
}
func StatSubComment(stat model.Stat, cancestat bool) (err error, errcode int) {
	return nil, 0
}
func StatComment(stat model.Stat, cancestat bool) (err error, errcode int) {
	return nil, 0
}
