package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//用户注册服务
type UserRegisterService struct {
	UserName    string `form:"username" json:"username" binding:"required,min=2,max=20"`
	Account     string `form:"account" json:"account" binding:"required,min=7,max=20"`
	Password    string `form:"password" json:"password" binding:"required,min=7,max=20"`
	RepPassword string `form:"reppassword" json:"reppassword" binding:"required,min=7,max=20"`
}

//用户登录服务
type UserLoginService struct {
	Account  string `form:"account" json:"account" binding:"required,min=7,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=7,max=20"`
}

//用户点赞列表
type UserStatListService struct {
	StatType model.StatType `form:"StatType" json:"StatType" binding:"omitempty,oneof=1 2"`
	UserId   uint           `form:"AuthorId" json:"AuthorId" binding:"omitempty"` //若为空，即为自己的点赞文章列表
	Paginationservice
}
type UpdateUsernameService struct {
	UserName string `form:"username" json:"username" binding:"required,min=2,max=20"`
}

func UserStatList() serializer.Response {
	return serializer.BuildResponse("xx")
}
func (service *UserRegisterService) Register() serializer.Response {
	if service.Password != service.RepPassword {
		return serializer.Err(serializer.ParamErr, errors.New("两次输入密码不同"))
	}
	count := 0

	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return serializer.Err(serializer.ParamErr, errors.New("用户名已存在"))
	}
	count = 0
	model.DB.Model(&model.User{}).Where("account = ?", service.Account).Count(&count)
	if count > 0 {
		return serializer.Err(serializer.ParamErr, errors.New("账号已存在"))
	}

	user := model.User{
		UserName: service.UserName,
		Account:  service.Account,
	}
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Err(serializer.ParamErr, errors.New("密码加密失败"))
	}

	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, errors.New("注册失败"))
	}
	return serializer.BuildResponse("注册成功")
}
func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	var user model.User

	if err := model.DB.Where("account = ?", service.Account).First(&user).Error; err != nil {
		return serializer.Err(serializer.ParamErr, errors.New("账户或密码错误"))
	}
	if user.CheckPassword(service.Password) == false {
		return serializer.Err(serializer.ParamErr, errors.New("账户或密码错误"))

	}
	s := sessions.Default(c)
	s.Clear()
	s.Options(sessions.Options{
		MaxAge: 86400,
	})
	s.Set("userID", user.ID) //gorm的自增ID
	s.Save()
	return serializer.BuildResponse("登录成功")
}
func (service *UpdateUsernameService) UpdateUsername(c *gin.Context) serializer.Response {
	user := model.GetcurrentUser(c)
	if user == nil {
		return serializer.Err(serializer.NoErr, errors.New("用户不存在"))
	}
	count := 0

	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return serializer.Err(serializer.ParamErr, errors.New("用户名已存在"))
	}
	if err := model.DB.Model(&user).Where("id=?", user.ID).Update("user_name", service.UserName).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, errors.New("修改失败"))

	}
	go func() {
		model.RedisWriteDB.Del(redis.UserIdKey(user.ID))
	}()
	return serializer.BuildResponse("修改成功")
}
