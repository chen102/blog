package model

type Article struct {
	ArticleId  uint
	AuthorId   uint
	Title      string `mapstructure:"title"`
	Stat       int    `mapstructure:"stat"`
	Time       string `mapstructure:"time"`
	Content    string `mapstructure:"content"`
	Tags       string `mapstructure:"tags"`
	Comment    []Comment
	CommentNum uint
}
type Comment struct {
	CommentId uint
	UserId    uint
	AuthorId  uint
	Time      string
	Content   string
	Stat      uint
}
