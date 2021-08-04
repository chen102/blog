package serializer

import (
	"blog/model"
)

type Article struct {
	AuthorId  string `json:"AuthorId"` //这里string的原因是model.Article是string
	ArticleId uint   `json:"ArticleId"`
	Title     string `json:"ArticleTitle"`
	Time      string `json:"ArticleTime"`
	Content   string `json:"ArticleContent"`
	Tags      string `json:"ArticleTags"`
}

func BuildArticle(article model.Article, id uint) Article {
	return Article{
		AuthorId:  article.AuthorId,
		ArticleId: id,
		Title:     article.Title,
		Time:      article.Time,
		Content:   article.Content,
		Tags:      article.Tags,
	}
}
func BuildArticleResponse(article model.Article, id uint) Response {
	return Response{
		Data: BuildArticle(article, id),
	}
}
