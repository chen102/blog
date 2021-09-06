package service

import (
	"blog/model"
	//"blog/model/redis"
	"blog/serializer"
	"blog/tool"
	"errors"
	"github.com/gin-gonic/gin"
	//"strconv"
)

//添加文章服务
type ArticleAddSservice struct {
	Title   string   `form:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content string   `form:"ArticleContent" json:"ArticleContent" binding:"required"`
	Tags    []string `form:"Tags" json:"Tags" binding:"omitempty"`
}

func (service *ArticleAddSservice) AddArticle(c *gin.Context) serializer.Response {
	user := model.GetcurrentUser(c)
	if user == nil {

		return serializer.Err(serializer.NoErr, errors.New("用户不存在"))
	}
	article := model.Article{
		Title:   service.Title,
		UserID:  user.ID,
		Content: service.Content,
		Tags:    tool.SliceToString(service.Tags),
	}
	if err := model.DB.Create(&article).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, err)
	}
	return serializer.BuildResponse("发表成功！")
}
