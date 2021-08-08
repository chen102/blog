package model

type Article struct {
	ArticleId uint
	AuthorId  uint
	Title     string `mapstructure:"title"`
	Time      string `mapstructure:"time"`
	Content   string `mapstructure:"content"`
	Tags      string `mapstructure:tags`
}
