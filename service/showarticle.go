package service

import (
	"blog/model"
	"blog/redis"
	"blog/serializer"
	"blog/tool"
	"strconv"

	//"blog/tool"
	//"errors"
	"github.com/gin-gonic/gin"
	//"encoding/json"
	//"fmt"
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
		user.ID = service.AuthorId
	} else { //若没有，默认是自己
		u := model.GetcurrentUser(c)
		if u != nil {
			user = *u
		}
	}
	if service.Count == 0 {
		service.Count = 5
	}
	data, err := redis.ShowArticleListCache(user.ID, service.Offset, service.Count, service.Type)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.RedisErr, err)
	} else if err == model.RedisNil {
		if err := model.DB.Where("user_id=?", user.ID).Find(&user.Articles).Error; err != nil { //直接查该用户所有文章写入redis，下次翻页，排序，都是在redis读服务器进行
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteArticleListCach(user.ID, user.Articles); err != nil {
			return serializer.Err(serializer.RedisErr, err)
		}

		return serializer.BuildArticleListResponse(user.Articles)
	}
	//手动处理data
	////Sort返回的结果为[]string，将string转为多个文章模型进行响应
	article := make([]model.Article, len(data)/5) //5个string为一个article,分别是id,title,time,stat,tags
	id := 0
	for i := 0; i < len(data); i++ {
		if i != 0 && i%5 == 0 {
			id++
		}
		//这样写真的很蠢
		switch i % 5 {
		case 0:
			articleid, err := strconv.Atoi(data[i])
			if err != nil {
				return serializer.Err(serializer.StrconvErr, err)
			}
			article[id].ID = uint(articleid)
		case 1:
			article[id].Title = data[i]
		case 2:
			article[id].UpdatedAt = tool.StringToTime(data[i])
		case 3:
			stat, err := strconv.Atoi(data[i])
			if err != nil {
				return serializer.Err(serializer.StrconvErr, err)
			}
			article[id].Stat = uint(stat)
		case 4:
			if data[i] != "" {

				article[id].Tags = data[i]
			}
		}
	}
	return serializer.BuildArticleListResponse(article)
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

		if err := mapstructure.WeakDecode(data, &article); err != nil {
			return serializer.Err(serializer.MapStructErr, err)
		}
	}
	//根据文章的用户ID查询该用户的名称
	var user []model.User
	username, err := redis.ShowUserNameCache(article.UserID)
	if err != nil && err != model.RedisNil {
		return serializer.Err(serializer.MysqlErr, err)
	} else if err == model.RedisNil { //若缓存未命中,更新缓存
		if err := model.DB.Select("user_name").First(&user, article.UserID).Error; err != nil {
			return serializer.Err(serializer.MysqlErr, err)
		}
		if err := redis.WriteUserNameCache(article.UserID, user[0].UserName); err != nil {

			return serializer.Err(serializer.RedisErr, err)
		}

		article.UserName = user[0].UserName
	}
	if username != "" {

		article.UserName = username
	}

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
