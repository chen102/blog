package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"github.com/gin-gonic/gin"
	//"strconv"
)

//动态服务
type DynamicService struct {
	Paginationservice
}

func (service *DynamicService) Dynamic(c *gin.Context) serializer.Response {
	me := model.GetcurrentUser(c)
	if service.Count == 0 {
		service.Count = 10
	}
	res, err := redis.ShowDynamicCache(me.ID, service.Offset, service.Count) //输出文章
	if err != nil && err != model.RedisNil {

		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil {
		var ids []string
		if ids, err = redis.ShowUserFollowerID(me.ID); err != nil && err != model.RedisNil {
			return serializer.Err(serializer.RedisErr, err)
		} else if err == model.RedisNil { //二级缓存失效,刷新关注用户列表
			users, err := db.UserFollowerList(me.ID, false)
			if err != nil {
				return serializer.Err(serializer.MysqlErr, err)
			}
			if err := redis.WriteFollowerListCache(me.ID, users, false); err != nil {
				return serializer.Err(serializer.RedisErr, err)
			}
			ids, _ = redis.ShowUserFollowerID(me.ID)
		}
		if len(ids) == 0 {
			return serializer.BuildResponse("没有关注用户")
		}
		var articles []model.Article
		if err := model.DB.Where("user_id IN (?)", ids).Order("ID desc").Find(&articles).Error; err != nil { //这里有个GORM的坑，要使用"(?),官方文档是?"

			return serializer.Err(serializer.MysqlErr, err)
		}
		articles, err = db.UserArticlesList(articles)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteDynamicCache(me.ID, articles); err != nil { //将关注用户文章写入缓存

			return serializer.Err(serializer.RedisErr, err)
		}
		if len(articles) > int(service.Count) {

			return serializer.BuildArticleListResponse(articles[:int(service.Count)])
		}
		return serializer.BuildArticleListResponse(articles)
	}

	return serializer.BuildArticleListResponse(res)
}
