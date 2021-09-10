package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"blog/tool"
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
		var useridsCache []string //从缓存拿用户关注ID
		var userids []int64       //从数据库拿用户ID
		if useridsCache, err = redis.ShowUserFollowerID(me.ID); err != nil && err != model.RedisNil {
			return serializer.Err(serializer.RedisErr, err)
		} else if err == model.RedisNil { //二级缓存失效,刷新关注用户列表
			//拿关注用户id
			userids, err = db.UserFollowerId(me.ID, false)
			if err != nil {
				return serializer.Err(serializer.MysqlErr, err)
			} else if len(userids) == 0 {

				return serializer.BuildResponse("没有关注用户")
			}
			//筛选id
			userids, err = redis.UserIncrementCache(userids)
			if err != nil {
				return serializer.Err(serializer.RedisErr, err)
			}
			//拿用户信息
			users, err := db.UserFollowerList(userids)
			if err != nil {
				return serializer.Err(serializer.MysqlErr, err)
			}
			if err := redis.WriteFollowerListCache(me.ID, users, false); err != nil {
				return serializer.Err(serializer.RedisErr, err)
			}
		}
		var ids []int64
		if len(useridsCache) != 0 {
			ids = tool.StringSliceTOIntSlice(useridsCache)
		} else {
			ids = userids
		}
		//拿到关注用户文章ID
		articleids, err := db.UserFollowerArticleID(ids)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		//筛选id
		articleids, err = redis.ArticleIncrementCache(me.ID, articleids)
		if err != nil {
			return serializer.Err(serializer.RedisErr, err)

		}
		//拿文章信息
		articles, err := db.UserArticlesList(articleids)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteDynamicCache(me.ID, articles); err != nil { //将关注用户文章写入缓存

			return serializer.Err(serializer.RedisErr, err)
		}
		//if len(articles) > int(service.Count) {

		//return serializer.BuildArticleListResponse(articles[:int(service.Count)])
		//}
		//return serializer.BuildArticleListResponse(articles)
		return serializer.BuildResponse("请再次执行")
	}

	return serializer.BuildArticleListResponse(res)
}
