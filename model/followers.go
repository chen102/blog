package model

type Follower struct {
	ID         uint `gorm:"primary_key"`
	UserID     uint
	FollowerID uint
	State      bool //取消关注
}
