package service

import (
	"blog/model/redis"
	"blog/serializer"
	"strconv"
)

//评论点赞服务
type StatCommentservice struct {
	UserId uint `form:"UserId" json:"UserId" binding:"required"`
	ArticleSservice
	CommentId uint `form:"CommentId json:"CommentId binding:required`
}

//文章点赞服务
type StatArticleservice struct {
	UserId uint `form:"UserId" json:"UserId" binding:"required"`
	ArticleSservice
}

func (service *StatCommentservice) StatComment() serializer.Response {
	if err := redis.StatComment(service.UserId, service.AuthorId, service.ArticleId, service.CommentId); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("stat ok!")
}

//文章点赞功能，因为文章是用hash存的，所以可以直接将stat作为文章的一个filed，然后通过这个字段排序，可以优先得到点赞高的文章
func (service *StatArticleservice) StatArticle() serializer.Response {
	if err := redis.StatArticle(service.UserId, service.AuthorId, service.ArticleId); err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("stat ok!")
}

//添加文章服务
type ArticleAddSservice struct {
	AuthorId uint     `form:"AuthorId" json:"AuthorId" binding:"required"`
	Title    string   `form:"ArticleTitle" json:"ArticleTitle" binding:"required,max=20"`
	Content  string   `form:"ArticleContent" json:"ArticleContent" binding:"required"`
	Tags     []string `form:"Tags" json:"Tags" binding:"omitempty"`
}

//添加评论服务
type ArticleCommentservice struct {
	UserId uint `form:"UserId" json:"UserId" binding:"requiredV"`
	ArticleSservice
	Content string `form:"CommentContent" json:"CommentContent" binding:"required"`
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
}
func (service *ArticleAddSservice) AddArticle() serializer.Response {
	id, err := redis.AddArticle(service.AuthorId, service.Title, service.Content, service.Tags)
	if err != nil {
		return serializer.Err(serializer.RedisErr, err)
	}
	return serializer.BuildResponse("Article ID:" + strconv.Itoa(id) + " ADD Succ！")
}
