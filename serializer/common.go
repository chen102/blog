package serializer

import (
	"model"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error",omitempty`
}
type Article struct {
}

func Build(article model.Article) Article {
	return Article{}
}
