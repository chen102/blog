package model

type Article struct {
	Title   string `mapstructure:"title"`
	Time    string `mapstructure:"time"`
	Content string `mapstructure:"content"`
}
