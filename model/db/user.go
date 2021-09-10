package db

import (
	. "blog/model"
)

func GetUser(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

//若fans为true,则返回粉丝列表,否则返回关注列表
func UserFollowerList(ids []int64) ([]User, error) {
	users := make([]User, len(ids))
	for k, id := range ids {
		//每个关注用户的关注、粉丝数、姓名、简介
		var username []string
		if err := DB.Model(&User{}).Where("id=?", id).Pluck("user_name", &username).Error; err != nil {
			return nil, err
		}
		var followercount int64
		var fanscount int64
		if err := DB.Model(&Follower{}).Where("user_id=? AND state=?", id, 0).Count(&followercount).Error; err != nil {
			return nil, err
		}
		if err := DB.Model(&Follower{}).Where("follower_id=? AND state=?", id, 0).Count(&fanscount).Error; err != nil {
			return nil, err
		}
		users[k].FansNum = uint(fanscount)
		users[k].FollowerNum = uint(followercount)
		users[k].ID = uint(id)
		users[k].UserName = username[0]
	}
	return users, nil
}
func UserFollowerId(userid uint, fans bool) ([]int64, error) {

	var ids []int64
	if !fans {
		if err := DB.Model(&Follower{}).Where("user_id=? AND state=?", userid, 0).Pluck("follower_id", &ids).Error; err != nil {
			return nil, err
		}

	} else {
		if err := DB.Model(&Follower{}).Where("follower_id=? AND state=?", userid, 0).Pluck("user_id", &ids).Error; err != nil {
			return nil, err
		}
	}
	return ids, nil
}
func UserArticleID(userid uint) ([]int64, error) {
	var articleids []int64
	if err := DB.Model(&Article{}).Where("user_id=?", userid).Pluck("id", &articleids).Error; err != nil {
		return nil, err
	}
	return articleids, nil
}
func UserFollowerArticleID(userids []int64) ([]int64, error) {
	var articleids []int64
	if err := DB.Model(&Article{}).Where("user_id IN (?)", userids).Pluck("id", &articleids).Error; err != nil { //这里有个GORM的坑，要使用"(?),官方文档是?"
		return nil, err

	}
	return articleids, nil
}

//给每篇文章点赞数,和用户名
func UserArticlesList(articleids []int64) ([]Article, error) {
	var articles []Article
	if err := DB.Where("id IN (?)", articleids).Order("ID desc").Find(&articles).Error; err != nil {
		return nil, err
	}
	for k, article := range articles {
		var username []string
		if err := DB.Model(&User{}).Where("id=?", article.UserID).Pluck("user_name", &username).Error; err != nil {
			return nil, err
		}
		articles[k].UserName = username[0]
		var count int64
		if err := DB.Model(&Stat{}).Where("type=? AND stat_id=? AND state=?", 0, article.ID, 0).Count(&count).Error; err != nil {

			return nil, err
		}
		articles[k].Stat = uint(count)
	}
	return articles, nil
}
func ExistUser(userid uint) bool {

	count := 0
	DB.Model(&User{}).Where("id = ?", userid).Count(&count)
	if count == 0 {
		return false
	}
	return true
}
