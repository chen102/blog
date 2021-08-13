package serializer

import (
	"blog/model"
	"strconv"
)

type Article struct {
	AuthorId   uint            `json:"AuthorId,omitempty"`
	ArticleId  uint            `json:"ArticleId,omitempty"`
	Title      string          `json:"ArticleTitle,omitempty"`
	Time       string          `json:"ArticleTime,omitempty"`
	Content    string          `json:"ArticleContent,omitempty"`
	Tags       string          `json:"ArticleTags,omitempty"`
	Comment    []model.Comment `json:"ArticleComment,omitempty"`
	CommentNum uint            `json:"ArticleCommentNum,omitempty"`
}
type Comment struct {
	CommentId uint
	UserId    uint
	AuthorId  uint
	Time      string
	Content   string
	Stat      uint
}

func BuildCommentList(comments []model.Comment) []Comment {
	comm := make([]Comment, 0)
	for _, comment := range comments {
		comm = append(comm, Comment{
			CommentId: comment.CommentId,
			UserId:    comment.UserId,
			AuthorId:  comment.AuthorId,
			Time:      comment.Time,
			Content:   comment.Content,
			Stat:      comment.Stat,
		})
	}
	return comm
}

//返回一篇文章的详细详细
func BuildArticle(article model.Article) Article {
	return Article{
		AuthorId:   article.AuthorId,
		ArticleId:  article.ArticleId,
		Title:      article.Title,
		Time:       article.Time,
		Content:    article.Content,
		Tags:       article.Tags,
		Comment:    article.Comment,
		CommentNum: article.CommentNum,
	}
}

//返回多条文章简要信息
func BuildArticleList(articles []model.Article) []Article {
	art := make([]Article, 0)
	for _, article := range articles {
		arttemp := Article{
			ArticleId: article.ArticleId,
			Title:     article.Title,
			Time:      article.Time,
			Tags:      article.Tags,
		}
		art = append(art, arttemp)
	}
	return art
}
func BuildArticleResponse(article model.Article) Response {
	return Response{
		Data: BuildArticle(article),
		Msg:  "article ID:" + strconv.Itoa(int(article.ArticleId)) + " Context Display Succ!",
	}
}
func BuildArticleListResponse(articles []model.Article) Response {
	return Response{
		Data: BuildArticleList(articles),
		Msg:  strconv.Itoa(len(articles)) + " articles Display Succ!",
	}
}
func BuildCommentListResponse(comments []model.Comment) Response {
	return Response{
		Data: BuildCommentList(comments),
		Msg:  strconv.Itoa(len(comments)) + " comment Display Succ!",
	}
}
