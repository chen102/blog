package api

import (
	"blog/serializer"
	"blog/service"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}

// @Summary 用户文章列表
// @Description 显示用户文章列表
// @Produce json
// @Param AuthorId body uint false "用户id，若为空，即为当前登录用户"
// @Param Rank body bool false "排序：若为空按时间排序，若为true按发布时间排序"
// @Param Offset body uint false "列表偏移量"
// @Param Count body uint false "列表一页请求的个数"
// @Success 200 {object} serializer.Article
// @Router /api/v0/article/list [post]
func ArticleList(c *gin.Context) {
	var service service.ArticleListservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.ArticleList(c)
		c.JSON(200, res)

	}
}

// @Summary 增加文章
// @Description 用户增加文章
// @Produce json
// @Param ArticleTitle body string true "文章标题名，最大20"
// @Param ArticleContent body string true "文章正文"
// @Param Tags body []string false  "文章标签"
// @Success 200 {string} string "发表成功"
// @Router /api/v0/article/add [post]
func AddArticle(c *gin.Context) {
	var service service.ArticleAddSservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.AddArticle(c)
		c.JSON(200, res)

	}

}

// @Summary 查看文章
// @Description 显示文章详情
// @Produce json
// @Param ArticleId body uint true "文章id"
// @Success 200 {object} []serializer.Article
// @Router /api/v0/article/show [get]
func ShowArticle(c *gin.Context) {
	var service service.ArticleService
	res := service.ShowArticle(c.DefaultQuery("id", "0"))
	c.JSON(200, res)
}

// @Summary 点赞
// @Description 用户给文章或者评论点赞
// @Produce json
// @Param StatId body uint true "点赞id"
// @Param  StatType body bool false "点赞类型:若为空，对文章点赞，为true点赞评论"
// @Param CancelStat body bool false "取消点赞:若为空，点赞，为true取消点赞"
// @Success 200 {string} string "成功"
// @Router /api/v0/article/stat [post]
func Stat(c *gin.Context) {
	var service service.StatService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Stat(c)
		c.JSON(200, res)
	}
}

// @Summary 评论
// @Description 用户评论文章或者回复评论
// @Produce json
// @Param id body uint true "文章id"
// @Param  Rev body uint false "评论id，若为空即为对文章评论"
// @Param LandlordId body uint false "回复评论楼主id，若是对评论回复这个字段不能为空"
// @Param  Content body string false "评论正文"
// @Success 200 {string} string "评论成功"
// @Router /api/v0/article/comment [post]
func Comment(c *gin.Context) {
	var service service.CommentService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Comment(c)
		c.JSON(200, res)
	}

}

// @Summary 评论列表
// @Description 显示文章评论列表
// @Produce json
// @Param ArticleId body uint true "文章id"
// @Param  Rev body uint false "评论id，若为空即为对文章评论"
// @Param Offset body uint false "列表偏移量"
// @Param Count body uint false "列表一页请求的个数"
// @Success 200 {object} []serializer.Comment
// @Router /api/v0/article/comment [post]
func ShowCommentList(c *gin.Context) {
	var service service.ArticleCommentListservice
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.CommentList()
		c.JSON(200, res)
	}

}

// @Summary 删除文章/评论
// @Description 用户删除指定的文章
// @Produce json
// @Param DeleteId body uint true "删除id"
// @Param ArticleId body uint false  "文章id,若是删除评论，这个字段不能为空"
// @Param Type body bool false "若为空，删文章，为true删评论"
// @Success 200 {string} string "删除成功"
// @Router /api/v0/article/delete [post]
func Delete(c *gin.Context) {
	var service service.DeleteService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, serializer.Err(serializer.ParamErr, err))
	} else {
		res := service.Delete(c)
		c.JSON(200, res)
	}

}
