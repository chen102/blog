package serializer

import (
	"blog/model"
	"time"
)

type Article struct {
	Id      uint      `json:"ArticleId"`
	Title   string    `json:"ArticleTitle"`
	Time    time.Time `json:"ArticleTime"`
	Content string    `json:"ArticleContent"`
}

func BuildArticle(article model.Article, id uint) Article {
	return Article{
		Id:      id,
		Title:   article.Title,
		Time:    article.Time,
		Content: article.Content,
	}
}
func BuildArticleResponse(article model.Article, id uint) Response {
	return Response{
		Data: BuildArticle(article, id),
	}
}
