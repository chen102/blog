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
	"blog/model/db"
	"strconv"
)

//文章服务
type ArticleService struct {
	ArticleId uint `form:"ArticleId" json:"ArticleId" binding:"required"`
}

//分页服务
type Paginationservice struct {
	Offset uint `form:"Offset" json:"Offset" binding:"omitempty"`
	Count  uint `form:"Count" json:"Count" binding:"omitempty"`
}

//用户文章列表
type ArticleListservice struct {
	Type     bool `form:"Rank" json:"Rank" binding:"omitempty"`
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
		//以前是全量，缓存索引失效，重新建立全部缓存(读mysql)，效率较低
		//增量缓存:先查id，去看此id缓存是否存在，若存在则刷新过期时间，然后将此id置0,将不存在的缓存id,查询mysql写入缓存,但有个问题，mysql不能输出全部数据了，缓存失效后第一次查询，数据不完整

		//拿文章ID
		articleids, err := db.UserArticleID(user.ID)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		//筛选ID
		articleids, err = redis.ArticleIncrementCache(user.ID, articleids)
		if err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		//拿文章
		user.Articles, err = db.UserArticlesList(articleids)
		if err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
		//写缓存
		if err := redis.WriteArticleListCach(user.ID, user.Articles); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
		//if len(user.Articles) > int(service.Count) {

		//return serializer.BuildArticleListResponse(user.Articles[:service.Count])
		//}
		//return serializer.BuildArticleListResponse(user.Articles)
		return serializer.BuildResponse("请再次执行")
	}
	//}
	//手动处理data
	////Sort返回的结果为[]string，将string转为多个文章模型进行响应
	return serializer.BuildArticleListResponse(article)
}
func (service *ArticleService) ShowArticle(id string) serializer.Response {
	artid, _ := strconv.Atoi(id)
	service.ArticleId = uint(artid)
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
		if err := model.DB.Model(&model.Stat{}).Where("stat_id=? AND state=?", article.ID, 0).Count(&count).Error; err != nil {

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
