package redis

import (
	"blog/model"
	"blog/tool"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

func DelOtherFansList(follower model.Follower) error {
	if err := model.RedisWriteDB.Del(UserFolloweListrKey(follower.UserID, 1)).Err(); err != nil {
		return err
	}
	return nil
}
func ExistUserFollowerList(userid uint) error {
	ok, err := model.RedisReadDB.Exists(UserFolloweListrKey(userid, 0)).Result()
	if err != nil && err != model.RedisNil {
		return err
	} else if ok == 0 {
		return model.RedisNil
	}
	return nil
}
func ShowUserFollowerID(userid uint) ([]string, error) {
	if err := ExistUserFollowerList(userid); err != nil {
		return nil, err
	}
	ids, err := model.RedisReadDB.SMembers(UserFolloweListrKey(userid, 0)).Result()
	if err != nil {
		return nil, err
	}
	return ids, nil
}
func WriteFollowerCache(follower model.Follower) error {
	if !follower.State {
		if err := model.RedisWriteDB.SAdd(UserFolloweListrKey(follower.UserID, 0), follower.FollowerID).Err(); err != nil {
			return err
		}

	} else { //取关
		if err := model.RedisWriteDB.SRem(UserFolloweListrKey(follower.UserID, 0), follower.FollowerID).Err(); err != nil {
			return err
		}
	}
	//集合(索引)绝对不能刷新存活时间，若刷新了，索引指向的底层缓存可能失效！！
	return nil
}
func WriteFollowerListCache(userid uint, users []model.User, fans bool) error {
	transactional := func(tx *redis.Tx) error {
		for _, user := range users {
			//关注列表集合
			if err := tx.SAdd(UserFolloweListrKey(userid, 0), user.ID).Err(); err != nil {
				return err
			}
			tx.Expire(UserFolloweListrKey(userid, 0), 24*time.Hour)
			exist, err := tx.Exists(UserIdKey(user.ID)).Result()
			if err != nil {
				return err
			} else if exist == 1 {
				tx.Expire(UserIdKey(user.ID), 24*time.Hour)
				continue
			}
			//用户信息
			if err := tx.HMSet(UserIdKey(user.ID), model.StructToMap(user)).Err(); err != nil {
				return err
			}
			tx.Expire(UserIdKey(user.ID), 24*time.Hour) //24小时存活

		}
		return nil
	}
	transactional1 := func(tx *redis.Tx) error {
		for _, user := range users {
			//粉丝列表集合
			if err := tx.SAdd(UserFolloweListrKey(userid, 1), user.ID).Err(); err != nil {
				return err
			}
			tx.Expire(UserFolloweListrKey(userid, 1), 24*time.Hour)
			exist, err := tx.Exists(UserIdKey(user.ID)).Result()
			if err != nil {
				return err
			} else if exist == 1 {
				tx.Expire(UserIdKey(user.ID), 24*time.Hour)
				continue
			}
			//用户信息
			if err := tx.HMSet(UserIdKey(user.ID), model.StructToMap(user)).Err(); err != nil {
				return err
			}
			tx.Expire(UserIdKey(user.ID), 24*time.Hour) //24小时存活

		}
		return nil
	}
	if !fans { //关注列表
		if err := model.RedisWriteDB.Watch(transactional, UserFolloweListrKey(userid, 0)); err != nil {
			return err
		}

	} else { //粉丝列表
		if err := model.RedisWriteDB.Watch(transactional1, UserFolloweListrKey(userid, 1)); err != nil {
			return err
		}

	}
	return nil
}
func ShowFollowerListCache(userid, offset, count uint, fans bool) ([]model.User, error) {
	var ids []string
	if !fans { //关注列表
		ok, err := model.RedisReadDB.Exists(UserFolloweListrKey(userid, 0)).Result()
		if err != nil && err != model.RedisNil {
			return nil, err
		} else if ok == 0 {
			return nil, model.RedisNil
		}
		ids, err = model.RedisReadDB.SMembers(UserFolloweListrKey(userid, 0)).Result()
		if err != nil {
			return nil, err
		}
		model.RedisWriteDB.Expire(UserFolloweListrKey(userid, 0), 24*time.Hour)
	} else { //粉丝列表

		ok, err := model.RedisReadDB.Exists(UserFolloweListrKey(userid, 1)).Result()
		if err != nil && err != model.RedisNil {
			return nil, err
		} else if ok == 0 {
			return nil, model.RedisNil
		}
		ids, err = model.RedisReadDB.SMembers(UserFolloweListrKey(userid, 1)).Result()
		if err != nil {
			return nil, err
		}
		model.RedisWriteDB.Expire(UserFolloweListrKey(userid, 1), 24*time.Hour)

	}

	//set没有分页功能，手动实现下
	page := tool.Pagination(ids, int(offset), int(count))
	if page == nil {
		return nil, nil
	}
	users := make([]model.User, len(page))
	for k, id := range page {
		data, err := model.RedisReadDB.HGetAll(UserStringIdKey(id)).Result()
		if err != nil {
			return nil, err
		}
		if data != nil {
			log.Println(data)
			mapstructure.WeakDecode(data, &users[k])
		}
		model.RedisWriteDB.Expire(UserStringIdKey(id), 1*time.Hour)
	}
	return users, nil
}
