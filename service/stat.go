package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"

	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	//"github.com/mitchellh/mapstructure"
	"encoding/json"
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
//用户点赞，写入redis，redis定时写入mysql点赞表，mysql计算点赞数，写入redis，
func StatArticle(stat model.Stat, cancestat bool) (err error, errcode int) {
	//文章模型
	var article model.Article
	if err := model.DB.First(&article, stat.ArticleID).Error; err != nil {
		return err, serializer.MysqlErr
	}
	//这里应该用redis查，初始化的时候，将用户id和用户名存在一个kv中，以后设计
	var user []model.User
	if err := model.DB.Select("user_name").First(&user, article.UserID).Error; err != nil {
		return err, serializer.MysqlErr
	}
	article.UserName = user[0].UserName
	var article_map map[string]interface{}
	data, _ := json.Marshal(&article)
	json.Unmarshal(data, &article_map)
	fmt.Println(article_map)
	if err := redis.Stat(stat, article_map); err != nil {
		return err, serializer.RedisErr
	}
	//取消点赞
	if cancestat {

	}
	//防止重复点赞
	if err := model.DB.Create(&stat).Error; err != nil {
		return err, serializer.MysqlErr
	}
	return nil, 0
}
func StatSubComment(stat model.Stat, cancestat bool) (err error, errcode int) {
	return nil, 0
}
func StatComment(stat model.Stat, cancestat bool) (err error, errcode int) {
	return nil, 0
}
