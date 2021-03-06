package serializer

import (
	"blog/model"
	"strconv"
	"time"
)

type Article struct {
	AuthorId   uint      `json:"AuthorId,omitempty"`
	AuthorName string    `json:"AuthorName,omitempty"`
	ArticleId  uint      `json:"ArticleId,omitempty"`
	Title      string    `json:"ArticleTitle,omitempty"`
	Time       time.Time `json:"ArticleTime,omitempty"`
	Content    string    `json:"ArticleContent,omitempty"`
	Tags       string    `json:"ArticleTags,omitempty"`
	Stat       uint      `json:"ArticleStat,omitempty"`
}

//返回一篇文章的详细详细
func BuildArticle(article model.Article) Article {
	return Article{
		AuthorId:   article.UserID,
		AuthorName: article.UserName,
		ArticleId:  article.ID,
		Title:      article.Title,
		Time:       article.UpdatedAt,
		Content:    article.Content,
		Tags:       article.Tags,
		Stat:       article.Stat,
	}
}

//返回多条文章简要信息
func BuildArticleList(articles []model.Article) []Article {
	art := make([]Article, 0)
	for _, article := range articles {
		arttemp := Article{
			ArticleId:  article.ID,
			AuthorName: article.UserName,
			Title:      article.Title,
			Time:       article.UpdatedAt,
			Tags:       article.Tags,
			Stat:       article.Stat,
		}
		art = append(art, arttemp)
	}
	return art
}
func BuildArticleResponse(article model.Article) Response {
	return Response{
		Data: BuildArticle(article),
		Msg:  "article ID:" + strconv.Itoa(int(article.ID)) + " Context Display Succ!",
	}
}
func BuildArticleListResponse(articles []model.Article) Response {
	return Response{
		Data: BuildArticleList(articles),
		Msg:  strconv.Itoa(len(articles)) + " articles Display Succ!",
	}
}
