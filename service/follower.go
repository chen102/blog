package service

//用户关注服务
type FollowerUserService struct {
	UserID         uint `form:"UserId" json:"UserId" binding:"required"`
	CancelFollower bool `form:"CancelFollower" json:"CancelFollower" binding:"omitempty"`
}

//关注列表服务
type UserFollowerListService struct {
	UserID uint `form:"UserId" json:"UserId" binding:"omitempty"` //不指定就是自己
}

//粉丝列表服务
type UserFansListService struct {
	UserID uint `form:"UserId" json:"UserId" binding:"omitempty"` //不指定就是自己
}
