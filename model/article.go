package model

type Article struct {
	ArticleId uint
	AuthorId  uint
	Title     string `mapstructure:"title"`
	Time      string `mapstructure:"time"`
	Content   string `mapstructure:"content"`
	Tags      string `mapstructure:"tags"`
	Comment   []string
}
type Comment struct {
	CommentId uint
	Content   []string
	stat      uint
}
