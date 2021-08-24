package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"
	//"blog/tool"
	"errors"
	"github.com/gin-gonic/gin"
	//"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	//"strconv"
)

//文章服务
type ArticleSservice struct {
	//AuthorId  uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	ArticleId uint `form:"ArticleId" json:"ArticleId" binding:"required"`
}

//分页服务
type Paginationservice struct {
	Offset uint `form:"Offset" json:"Offset" binding:"omitempty"`
	Count  uint `form:"Count" json:"Count" binding:"omitempty"`
}

//文章分页服务
type ArticleListservice struct {
	Type     bool `form:"rank" json:"rank" binding:"omitempty"`
	AuthorId uint `form:"AuthorId" json:"AuthorId" binding:"omitempty"` //若为空，即为自己的文章列表
	Paginationservice
}

//评论分页服务,若评论超过5条，需发起请求获取全部评论(分页的形式)
type ArticleCommentListservice struct {
	ArticleSservice
	Paginationservice
}

func (service *ArticleCommentListservice) ArticleCommentList() serializer.Response {
	//var comments []model.Comment
	//res, err := redis.ShowAllComment(service.AuthorId, service.ArticleId, service.Offset, service.Count)
	//if err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	//for _, comment := range res {
	//var temp model.Comment
	//strconv.Unquote(comment)
	//if len(comment) == 0 {
	//continue
	//}
	//if err := json.Unmarshal([]byte(comment), &temp); err != nil {
	//return serializer.Err(serializer.StrconvErr, err)
	//}
	//stat, err := redis.GetStat(service.AuthorId, service.ArticleId, temp.CommentId)
	//if err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	//temp.Stat = stat
	//comments = append(comments, temp)

	//}
	//return serializer.BuildCommentListResponse(comments)
	return serializer.BuildResponse("xx")
}

func (service *ArticleListservice) ArticleList(c *gin.Context) serializer.Response {
	var user model.User
	if service.AuthorId != 0 { //指定了用户
		if err := model.DB.Where("user_id=?", service.AuthorId).First(&user).Error; err != nil {
			fmt.Println(service.AuthorId)
			return serializer.Err(serializer.MysqlErr, err)
		}
	} else {

		onse := model.GetcurrentUser(c) //没有指定用户默认就是自己的列表
		if onse == nil {

			return serializer.Err(serializer.NoErr, errors.New("用户不存在"))
		}
		user = *onse
	}
	if service.Count == 0 {
		service.Count = 5
	}
	if err := model.DB.Limit(service.Count).Offset(service.Offset).Where("user_id=?", user.ID).Find(&user.Articles).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, err)
	}
	return serializer.BuildArticleListResponse(user.Articles)
	//res, err := redis.ListArticle(service.AuthorId, service.Offset, service.Count, service.Type)
	//if err != nil {

	//return serializer.Err(serializer.RedisErr, err)
	//}
	////Sort返回的结果为string，将string转为多个文章模型进行响应
	//article := make([]model.Article, len(res)/4) //4个string为一个article,分别是id,title,time,tags
	//id := 0
	//for i := 0; i < len(res); i++ {
	//if i != 0 && i%4 == 0 {
	//id++
	//}
	//switch i % 4 {
	//case 0:
	//articleid, err := strconv.Atoi(res[i])
	//if err != nil {
	//return serializer.Err(serializer.StrconvErr, err)
	//}
	//article[id].ArticleId = uint(articleid)
	//case 1:
	//article[id].Title = res[i]
	//case 2:
	//article[id].Time = res[i]
	//case 3:
	//if res[i] != "" {

	//article[id].Tags = res[i]
	//}
	//}
	//}
	//return serializer.BuildArticleListResponse(article)
	//return serializer.BuildResponse("xx")
}
func (service *ArticleSservice) DeleteArticle() serializer.Response {
	return serializer.Response{}

}
func (service *ArticleSservice) ShowArticle() serializer.Response {
	var article model.Article
	data, err := redis.ShowArticleCache(service.ArticleId)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil { //若缓存未命中,更新缓存
		if err := model.DB.First(&article, service.ArticleId).Error; err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteArticleCache(model.StructToMap(article)); err != nil {

			return serializer.Err(serializer.RedisErr, err)
		}

	}
	if data != nil {

		if err := mapstructure.Decode(data, &article); err != nil {
			return serializer.Err(serializer.MapStructErr, err)
		}
	}
	//根据文章的用户ID查询该用户的名称
	var user []model.User
	if err := model.DB.Select("user_name").First(&user, article.UserID).Error; err != nil {
		return serializer.Err(serializer.MysqlErr, err)
	}
	article.UserName = user[0].UserName

	return serializer.BuildArticleResponse(article)
	//var article model.Article
	//data, err := redis.ShowArticle(service.AuthorId, service.ArticleId)
	//if err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	//if err := mapstructure.Decode(data, &article); err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	//article.ArticleId = service.ArticleId
	//article.AuthorId = service.AuthorId
	////显示评论
	//var comments []model.Comment

	//comment, err := redis.ShowComment(service.AuthorId, service.ArticleId)
	//if err != nil {

	//return serializer.Err(serializer.RedisErr, err)
	//}
	//commentnumString := comment[len(comment)-1] //查询的评论数
	//commentnumInt, err := strconv.Atoi(commentnumString)
	//if err != nil {
	//return serializer.Err(serializer.StrconvErr, err)
	//}
	//for _, comment := range comment[:len(comment)-1] {
	//var temp model.Comment
	//strconv.Unquote(comment)
	//if err := json.Unmarshal([]byte(comment), &temp); err != nil {
	//return serializer.Err(serializer.StrconvErr, err)
	//}
	//stat, err := redis.GetStat(service.AuthorId, service.ArticleId, temp.CommentId)
	//if err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	//temp.Stat = stat
	//comments = append(comments, temp)

	//}
	//article.CommentNum = uint(commentnumInt)
	//article.Comment = comments
	//return serializer.BuildArticleResponse(article)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
