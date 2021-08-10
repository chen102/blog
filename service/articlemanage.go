package service

import (
	"blog/model"
	"blog/model/redis"
	"blog/serializer"
	//"blog/tool"
	//"encoding/json"
	//"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type ArticleAddSservice struct {
	AuthorId uint     `form:"AuthorId" json:"AuthorId" binding:"required"`
	Title    string   `form:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content  string   `form:"ArticleContent" json:"ArticleContent" binding:"required"`
	Tags     []string `form:"Tags" json:"Tags" binding:"omitempty"`
}
type ArticleSservice struct {
	AuthorId  uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	ArticleId uint `form:"ArticleId" json:"ArticleId" binding:"required"`
}
type ArticleListservice struct {
	AuthorId uint `form:"AuthorId" json:"AuthorId" binding:"required"`
	Offset   uint `form:"ArticleOffset" json:"ArticleOffset" binding:"omitempty"`
	Count    uint `form:"ArticleCount" json:"ArticleCount" binding:"omitempty"`
}
type ArticleCommentservice struct {
	UserId    uint   `form:"UserId" json:"UserId" binding:"required"`
	AuthorId  uint   `form:"AuthorId" json:"AuthorId" binding:"required"`
	ArticleId uint   `form:"ArticleId" json:"ArticleId" binding:"required"`
	Content   string `form:"CommentContent" json:"CommentContent" binding:"required"`
	//子评论，回复评论
	CommentId uint `form:"CommentId" json:"CommentId" binding:"omitempty"`
	ToUser    uint `form:"ToUser" json:"ToUser" binding:"omitempty"`
}
type StatCommentservice struct {
}
type StatArticleservice struct {
}

func (service *StatCommentservice) StatComment() serializer.Response {
	return serializer.BuildResponse("xx")
}
func (service *StatArticleservice) StatArticle() serializer.Response {
	return serializer.BuildResponse("xx")
}

//评论功能设计
//有序集合存储评论ID，根据评论点赞排行
//哈希存储评论和子评论
//对文章评论：生成新评论ID，加入评论集，写入评论
//对评论评论：找到评论ID，写入评论评论
func (service *ArticleCommentservice) ArticleComment() serializer.Response {
	id, err := redis.AddComment(service.UserId, service.AuthorId, service.ArticleId, service.Content)
	if err != nil {

		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("Comment ID:" + strconv.Itoa(int(id)) + " ADD Succ！")
	//if service.CommentId != 0 { //子评论
	//if service.ToUser == 0 {
	//return serializer.BuildResponse("ToUser is Null! ")
	//}
	//comment["touser"] = service.ToUser
	//comment_json, err := json.Marshal(comment)
	//if err != nil {
	//return serializer.Err(serializer.StrconvErr, err)
	//}
	//ok, err := model.Redisdb.HExists(model.ArticleCommentIDKey(service.CommentId), "comment").Result()
	//if !ok {
	//return serializer.BuildResponse("Comment ID:" + strconv.Itoa(int(service.CommentId)) + " NOT EXIST")
	//}
	//if err = model.Redisdb.HSet(model.ArticleCommentIDKey(service.CommentId), "subcomments", comment_json).Err(); err != nil {

	//return serializer.Err(serializer.RedisErr, err)
	//}
	//return serializer.BuildResponse("Recv Succ!")
	//} else { //评论
	////开始事务
	////生成评论id,保证id并发安全
	//pipe := model.Redisdb.TxPipeline()
	//if err := pipe.SetNX(model.GetArticleCommentsIDKey(service.ArticleId), 1, 0).Err(); err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}
	////wahch用法不对
	////if err := pipe.Watch(model.GetArticleCommentsIDKey(service.ArticleId)).Err(); err != nil {
	////return serializer.Err(serializer.RedisErr, err)
	////}
	//commentid, err := pipe.Get(model.GetArticleCommentsIDKey(service.ArticleId)).Result()
	//if err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}

	//pipe.Incr(model.GetArticleCommentsIDKey(service.ArticleId))
	////将评论id写入有序集合
	//if err = pipe.ZAdd(model.ArticleCommentRankKey(service.ArticleId), model.InitCommentRank(commentid)).Err(); err != nil {

	//return serializer.Err(serializer.RedisErr, err)
	//}
	//if _, err := pipe.Exec(); err != nil {
	//return serializer.Err(serializer.RedisErr, err)
	//}

	//}
}
func (service *ArticleListservice) ArticleList() serializer.Response {
	res, err := redis.ListArticle(service.AuthorId, service.Offset, service.Count)
	if err != nil {

		return serializer.Err(serializer.RedisErr, err)
	}
	//Sort返回的结果为string，将string转为多个文章模型进行响应
	article := make([]model.Article, len(res)/4) //4个string为一个article,分别是id,title,time,tags
	id := 0
	for i := 0; i < len(res); i++ {
		if i != 0 && i%4 == 0 {
			id++
		}
		switch i % 4 {
		case 0:
			articleid, err := strconv.Atoi(res[i])
			if err != nil {
				return serializer.Err(serializer.StrconvErr, err)
			}
			article[id].ArticleId = uint(articleid)
		case 1:
			article[id].Title = res[i]
		case 2:
			article[id].Time = res[i]
		case 3:
			if res[i] != "" {

				article[id].Tags = res[i]
			}
		}
	}
	return serializer.BuildArticleListResponse(article)
}
func (service *ArticleAddSservice) AddArticle() serializer.Response {
	id, err := redis.AddArticle(service.AuthorId, service.Title, service.Content, service.Tags)
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("Article ID:" + strconv.Itoa(id) + " ADD Succ！")
}
func (service *ArticleSservice) DeleteArticle() serializer.Response {
	return serializer.Response{}

}
func (service *ArticleSservice) ShowArticle() serializer.Response {

	var article model.Article
	data, err := redis.ShowArticle(service.AuthorId, service.ArticleId)
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	if err := mapstructure.Decode(data, &article); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	article.ArticleId = service.ArticleId
	article.AuthorId = service.AuthorId
	//显示评论
	comment, err := redis.ShowComment(service.AuthorId, service.ArticleId)
	if err != nil {

		return serializer.Err(serializer.RedisErr, err)
	}
	article.Comment = comment
	return serializer.BuildArticleResponse(article)
}
func (service *ArticleSservice) UpdateArticle() serializer.Response {
	return serializer.Response{}
}
