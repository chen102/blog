package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"
	//"blog/tool"
	//"strconv"
	//"blog/tool"
	//"errors"
	"github.com/gin-gonic/gin"
	//"encoding/json"
	//"fmt"
	//"github.com/mitchellh/mapstructure"
	//"strconv"
	"blog/model/db"
)

//文章服务
type ArticleService struct {
	//AuthorId  uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	ArticleId uint `form:"ArticleId" json:"ArticleId" binding:"required"`
}

//分页服务
type Paginationservice struct {
	Offset uint `form:"Offset" json:"Offset" binding:"omitempty"`
	Count  uint `form:"Count" json:"Count" binding:"omitempty"`
}

//用户文章列表
type ArticleListservice struct {
	Type     bool `form:"rank" json:"rank" binding:"omitempty"`
	AuthorId uint `form:"AuthorId" json:"AuthorId" binding:"omitempty"` //若为空，即为自己的文章列表
	Paginationservice
}

func (service *ArticleListservice) ArticleList(c *gin.Context) serializer.Response {
	var user model.User
	if service.AuthorId != 0 { //指定了用户
		user.ID = service.AuthorId
	} else { //若没有，默认是自己
		u := model.GetcurrentUser(c)
		if u != nil {
			user = *u
		}
	}
	if service.Count == 0 {
		service.Count = 5
	}
	article, err := redis.ShowArticleListCache(user.ID, service.Offset, service.Count, service.Type)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil {
		if err := model.DB.Where("user_id=?", user.ID).Order("ID desc").Find(&user.Articles).Error; err != nil { //直接查该用户所有文章写入redis，下次翻页，排序，都是在redis读服务器进行
			return serializer.Err(serializer.MysqlErr, err)
		}
		user.Articles, err = db.UserArticlesList(user.Articles)
		if err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteArticleListCach(user.ID, user.Articles); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		if len(user.Articles) > int(service.Count) {

			return serializer.BuildArticleListResponse(user.Articles[:service.Count])
		}
		return serializer.BuildArticleListResponse(user.Articles)
	}
	//}
	//手动处理data
	////Sort返回的结果为[]string，将string转为多个文章模型进行响应
	return serializer.BuildArticleListResponse(article)
}
func (service *ArticleService) DeleteArticle() serializer.Response {
	return serializer.Response{}

}
func (service *ArticleService) ShowArticle() serializer.Response {
	if !db.ExistArticle(service.ArticleId) {
		return serializer.BuildResponse("没有此文章")
	}
	var article model.Article
	data, err := redis.ShowArticleCache(service.ArticleId)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil { //若缓存未命中,更新缓存
		if err := model.DB.First(&article, service.ArticleId).Error; err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		var user []model.User
		//获取用户名
		if err := model.DB.Select("user_name").First(&user, article.UserID).Error; err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		article.UserName = user[0].UserName

		//获取点赞数
		var count int64
		if err := model.DB.Model(&model.Stat{}).Where("article_id=? AND Stat=?", article.ID, 0).Count(&count).Error; err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
		article.Stat = uint(count)
		if err := redis.WriteArticleCache(model.StructToMap(article)); err != nil {

			return serializer.Err(serializer.RedisErr, err)
		}
		//获取评论

		return serializer.BuildArticleResponse(article)
		//return serializer.BuildCommentListResponse(comment)
	}
	return serializer.BuildArticleResponse(data[0])
}
func (service *ArticleService) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
