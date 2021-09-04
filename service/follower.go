package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

//用户关注服务
type FollowerUserService struct {
	UserId         uint `form:"UserId" json:"UserId" binding:"required"`
	CancelFollower bool `form:"CancelFollower" json:"CancelFollower" binding:"omitempty"`
}

//关注列表服务
type UserFollowerListService struct {
	UserId uint `form:"UserId" json:"UserId" binding:"omitempty"` //不指定就是自己
	//0 1 2 关注、粉丝、互关列表
	ListType uint `form:"Type" json:"Type" binding:"omitempty,oneof=1 2"`
	Paginationservice
}

func (service *FollowerUserService) FollowerUser(c *gin.Context) serializer.Response {
	count := 0
	model.DB.Model(&model.User{}).Where("id = ?", service.UserId).Count(&count)
	if count == 0 {
		return serializer.BuildResponse("没有此用户")
	}
	me := model.GetcurrentUser(c)
	follow := model.Follower{
		UserID:     me.ID,
		FollowerID: service.UserId,
	}
	//删除别人粉丝列表缓存(粉丝列表要考虑并发性),配合redis事务，可控制并发
	if err := redis.DelOtherFansList(follow); err != nil {

		return serializer.Err(serializer.RedisErr, err)
	}
	//写入mysql
	if !service.CancelFollower {

		if err := model.DB.Create(&follow).Error; err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
	} else {
		if err := model.DB.Model(&follow).Where("user_id=? AND follower_id=?", follow.UserID, follow.FollowerID).Update("stat", true).Error; err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
	}
	follow.Stat = service.CancelFollower
	//写入缓存
	if err := redis.ExistUserFollowerList(follow.UserID); err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == nil {
		if err := redis.WriteFollowerCache(follow); err != nil { //若缓存存在可直接同步修改缓存,自己的关注列表不用考虑并发性

			return serializer.Err(serializer.RedisErr, err)
		}
	} else if err == model.RedisNil { //缓存未命中，构建缓存
		//从关注表，获取关注用户信息写入cache
		users, err := db.UserFollowerList(me.ID, false)
		if err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteFollowerListCache(me.ID, users, false); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}
	}
	return serializer.BuildResponse("关注成功！")
}
func (service *UserFollowerListService) UserFollowerList(c *gin.Context) serializer.Response {
	if service.UserId == 0 {
		session := sessions.Default(c)
		service.UserId = session.Get("userID").(uint)
	}
	if service.Count == 0 {
		service.Count = 5
	}
	var Users []model.User
	switch service.ListType {
	case 0:
		res, err, errcode := followList(service.UserId, service.Offset, service.Count)
		if err != nil {
			return serializer.Err(errcode, err)
		}
		Users = res

	case 1:
		res, err, errcode := fansList(service.UserId, service.Offset, service.Count)
		if err != nil {
			return serializer.Err(errcode, err)
		}
		Users = res

	case 2:
		res, err, errcode := mutualfollowingList(service.UserId, service.Offset, service.Count)
		if err != nil {
			return serializer.Err(errcode, err)
		}
		Users = res
	}
	return serializer.BuildUserListResponse(Users)
}
func followList(userid, offset, count uint) ([]model.User, error, int) {
	log.Println("DEBUG1")
	res, err := redis.ShowFollowerListCache(userid, offset, count, false)
	if err != nil && err != model.RedisNil {
		return nil, err, serializer.RedisErr
	} else if err == model.RedisNil {
		users, err := db.UserFollowerList(userid, false)
		if err != nil {
			return nil, err, serializer.RedisErr
		}
		if err := redis.WriteFollowerListCache(userid, users, false); err != nil {
			return nil, err, serializer.RedisErr
		}
		return users, nil, 0

	}
	return res, nil, 0
}
func fansList(userid, offset, count uint) ([]model.User, error, int) {
	res, err := redis.ShowFollowerListCache(userid, offset, count, true)
	if err != nil && err != model.RedisNil {
		return nil, err, serializer.RedisErr
	} else if err == model.RedisNil {
		users, err := db.UserFollowerList(userid, true)
		if err != nil {
			return nil, err, serializer.MysqlErr
		}
		if err := redis.WriteFollowerListCache(userid, users, true); err != nil {
			return nil, err, serializer.RedisErr
		}
		return users, nil, 0

	}
	return res, nil, 0
}
func mutualfollowingList(follower, offset, count uint) ([]model.User, error, int) {
	return nil, nil, 0
}
