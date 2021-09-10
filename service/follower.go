package service

import (
	"blog/model"
	"blog/model/db"
	"blog/redis"
	"blog/serializer"
	"github.com/gin-gonic/gin"
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

	if !db.ExistUser(service.UserId) {

		return serializer.BuildResponse("没有此用户")
	}
	me := model.GetcurrentUser(c)
	follow := model.Follower{
		UserID:     me.ID,
		FollowerID: service.UserId,
	}
	//缓存策略,先更db后删缓存
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
	//删除别人粉丝列表缓存(粉丝列表要考虑并发性:a查看自己的粉丝列表，同时b关注了a，缓存无数据，a从数据库读到数据后，准备回写缓存，此刻b正好要更新数据库，更新完后，又去把缓存更新了，那请求a再往缓存中写的就是旧数据，属于脏数据),配合redis事务，可控制并发
	if err := redis.DelOtherFansList(follow); err != nil {

		return serializer.Err(serializer.RedisErr, err)
	}
	follow.State = service.CancelFollower
	//写入缓存
	if err := redis.ExistUserFollowerList(follow.UserID); err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == nil {
		if err := redis.WriteFollowerCache(follow); err != nil { //若缓存存在可直接同步修改缓存,自己的关注列表不用考虑并发性:我不能同时关注别人和看自己的粉丝列表

			return serializer.Err(serializer.RedisErr, err)
		}
	} else if err == model.RedisNil { //缓存未命中，构建缓存
		//从关注表，获取关注用户信息写入cache
		userids, err := db.UserFollowerId(me.ID, false)
		if err != nil {

			return serializer.Err(serializer.MysqlErr, err)
		}
		users, err := db.UserFollowerList(userids)
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
		me := model.GetcurrentUser(c)
		service.UserId = me.ID
	} else {
		if !db.ExistUser(service.UserId) {
			return serializer.BuildResponse("没有此用户")
		}
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
	res, err := redis.ShowFollowerListCache(userid, offset, count, false)
	if err != nil && err != model.RedisNil {
		return nil, err, serializer.RedisErr
	} else if err == model.RedisNil {
		userids, err := db.UserFollowerId(userid, false)
		if err != nil {

			return nil, err, serializer.MysqlErr
		}
		users, err := db.UserFollowerList(userids)
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
		userids, err := db.UserFollowerId(userid, true)
		if err != nil {

			return nil, err, serializer.MysqlErr
		}
		users, err := db.UserFollowerList(userids)
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
