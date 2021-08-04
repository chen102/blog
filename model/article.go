package model

type Article struct {
	AuthorId string `mapstructure:authorid` //这里是string而不是uint的原因是因为mapstructure中没有string转uint的方法
	Title    string `mapstructure:"title"`
	Time     string `mapstructure:"time"`
	Content  string `mapstructure:"content"`
	Tags     string `mapstructure:tags`
}
