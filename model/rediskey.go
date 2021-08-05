//Redis所有Key
//命名规范:   模块名:业务逻辑:value类型
package model

import (
	"blog/tool"
)

//文章全局id键 articlemanager:articleid:xx
func GetArticleIDKey() string {
	return "articlemanager:articleid:int"
}

//文章键 articlemanager:articleid:xx:hash
func ArticleIdKey(articleid string) string {
	return tool.StrSplicing("articlemanager:articleid:", articleid, ":hash")
}

//用户文章合集键 articlemanager:authorid:xx:articles:set
func AuthorArticlesKey(authorid string) string {
	return tool.StrSplicing("articlemanager:authorid:", authorid, ":articles:set")
}

//文章标签键 articlemanager:articleid:xx:tags:set
func ArticleTagsKey(articleid string) string {
	return tool.StrSplicing("articlemanager:articleid", articleid, ":tags:set") //article:articleID:tags
}

//标签键 articlemanager:articleid:xx:tags:set
//这里等于是把所有有xx标签的文章都找到了(不管是不是自己的)，若需按自己文章标签分类，与自己文章id取交集即可 tips:这里应该是有更好的方法实现按标签分类自己的文章(或许这个关系型，确实不应该用redis做)
func TagKey(tag string) string {
	return tool.StrSplicing("articlemanager:articleid", tag, ":tags:set") //article:articleID:tags
}
