package api

import (
	"blog/serializer"
	"blog/service"
	"github.com/gin-gonic/gin"
)

// @Summary 注册用户
// @Description 用户输入账号密码等信息进行注册
// @Produce json
// @Param username body string true "用户名,长度大于2，小于20"
// @Param account body string true "账号,长度大于7，小于20"
// @Param password body string true "密码,长度大于7，小于20"
// @Param reppassword body string true "重复密码"
// @Success 200 {string} string "注册成功"
// @Router /api/v0/user/register [post]
func UserRegister(c *gin.Context) {
	var service service.UserRegisterService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Register()
		c.JSON(200, res)

	}
}

// @Summary 登录用户
// @Description 用户输入账号密码进行登录
// @Produce json
// @Param account body string true "账号,长度大于7，小于20"
// @Param password body string true "密码,长度大于7，小于20"
// @Success 200 {string} string "登录成功"
// @Router /api/v0/user/login [post]
func UserLogin(c *gin.Context) {
	var service service.UserLoginService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Login(c)
		c.JSON(200, res)

	}
}

// @Summary 用户喜欢文章列表
// @Description 显示用户点赞的文章列表
// @Produce json
// @Param AuthorId body uint false "用户id，若为空，即为当前登录用户"
// @Param Offset body uint false "列表偏移量"
// @Param Count body uint false "列表一页请求的个数"
// @Success 200 {object} []serializer.Article
// @Router /api/v0/user/like [post]
func UserArticlesLike(c *gin.Context) {
	var service service.UserStatListService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.UserStatArticlesList(c)
		c.JSON(200, res)

	}

}

// @Summary 关注用户
// @Description 用户关注另一个用户
// @Produce json
// @Param UserId body uint true "用户id"
// @Param CancelFollower body bool false "若为true对该用户取关"
// @Success 200 {string} string "关注成功"
// @Router /api/v0/follower [post]
func UserFollowerUser(c *gin.Context) {
	var service service.FollowerUserService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.FollowerUser(c)
		c.JSON(200, res)

	}

}

// @Summary 用户关注/粉丝列表
// @Description 显示用户关注的用户列表信息
// @Produce json
// @Param UserId body uint false "用户id，若为空，即为当前登录用户"
// @Param Type body uint false "1:只能是1或为空，空:关注列表，1:"粉丝列表
// @Param Offset body uint false "列表偏移量"
// @Param Count body uint false "列表一页请求的个数"
// @Success 200 {object} []serializer.User
// @Router /api/v0/follower/list [post]
func UserFollowerList(c *gin.Context) {

	var service service.UserFollowerListService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.UserFollowerList(c)
		c.JSON(200, res)

	}

}

// @Summary 用户关注动态
// @Description 显示关注用户的动态
// @Produce json
// @Param Offset body uint false "列表偏移量"
// @Param Count body uint false "列表一页请求的个数"
// @Success 200 {object} []serializer.Article
// @Router /api/v0/follower/dynamic [post]
func UserDynamicList(c *gin.Context) {

	var service service.DynamicService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Dynamic(c)
		c.JSON(200, res)

	}

}
