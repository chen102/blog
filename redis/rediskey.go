//Redis所有Key
//命名规范:   模块名:业务逻辑:存储的东西:value类型
package redis

import (
	"blog/tool"
	"github.com/go-redis/redis"
	"strconv"
)

const RedisNil = redis.Nil

//redis事务逻辑
var Transactional *redis.Tx

func InitCommentRank(member interface{}) redis.Z {
	return redis.Z{Score: 0, Member: member}
}

//Sort命令的参数
func SortArgs(by string, offset, count int64, get []string, order string, alpha bool) *redis.Sort {
	return &redis.Sort{
		By:     by,
		Offset: offset,
		Count:  count,
		Get:    get,
		Order:  order,
		Alpha:  alpha,
	}
}

//获取key排序后指定的字段
func GetSort(key string, returnId bool, strs ...string) []string {
	get := make([]string, 0)
	if returnId {
		get = append(get, "#") //GET #返回自身
	}
	for _, str := range strs {
		get = append(get, tool.StrSplicing(key, "->", str))
	}
	return get
}
func UintToStr(str uint) string {
	return strconv.Itoa(int(str))
}

//文章点赞数
//statservice:articleid:xx:likes:int
func ArticleStatKey(artid uint) string {
	return tool.StrSplicing("statservice:articleid:", UintToStr(artid), ":likes:int")
}

//用户所有的点赞文章
func UserStatArticleKey(userid, artid uint) string {
	return tool.StrSplicing("userservice:userid:", UintToStr(userid), ":articleid:", UintToStr(artid), "article:hash")
}

//文章主键键 articlemanager:userid:xx:articlid:int
func GetArticleIDKey(userid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), "articlid:int")
}

//评论主键 articlemanager:userid:xx:articleid:xx:commentid:int
func GetArticleCommentsIDKey(userid, articleid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:int")
}

//子评论主键
func GetArticleSubCommentsIDKey(userid, articleid, commentid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:", UintToStr(commentid), ":subcommentid:int")
}

//文章键 articlemanager:userid:xx:articleid:xx:hash
//value:用户的文章内容
func ArticleIdKey(userid, articleid uint) string {

	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":hash")
}
func ArticleIdStringKey(userid, articleid string) string {

	return tool.StrSplicing("articlemanager:userid:", userid, ":articleid:", articleid, ":hash")
}

//用户文章合集键 articlemanager:userid:xx:articleid:set
//value:用户的所有文章id
func AuthorArticlesKey(userid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:set")
}

//文章标签键 articlemanager:userid:xx:articleid:xx:tag:set
//vale:用户每篇文章的标签
func ArticleTagsKey(userid, articleid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":tags:set") //article:articleID:tags
}
func ArticleTagsStringKey(userid, articleid string) string {
	return tool.StrSplicing("articlemanager:userid", userid, ":articleid:", articleid, ":tags:set") //article:articleID:tags
}

//标签键 articlemanager:userid:xx:tag:xx:articleid:set
//value用户的每个标签对应的文章id
func TagKey(userid uint, tag string) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":tag:", tag, ":articleid:set")
}

//文章评论合集 articlemanager:userid:xx:articleid:xx:commentid:zset
//value:用户文章的所有评论id
func ArticleCommentRankKey(userid, articleid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:zset")
}

//文章子评论合集 articlemanager:userid:xx:articleid:xx:commentid:xx:subcommentid:zset
//value:用户文章评论的每个子评论id
func ArticleSubCommentRankKey(userid, articleid, commentid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:", UintToStr(commentid), ":subcommentid:zset")
}

//文章评论键 articlemanager:userid:xx:articleid:xx:commentid:xx:hash
//value:子\评论内容
func ArticleCommentIDStringKey(userid, articlid, commentid string) string {
	return tool.StrSplicing("articlemanager:userid:", userid, ":articleid:", articlid, ":commentid:", commentid, ":hash")
}
func ArticleCommentIDKey(userid, articleid, commentid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:", UintToStr(commentid), ":hash")
}

//评论点赞数 articlemanager:userid:xx:articleid:xx:commentid:xx:stats:int
//value:每条评论的点赞数
func ArticleCommentStatKey(userid, articleid, commentid uint) string {
	return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":articleid:", UintToStr(articleid), ":commentid:", UintToStr(commentid), ":stats:int")
}

//用户点赞的文章/评论合集 articlemanager:userid:xx:comment:stat:set
//评论value:用户:文章:评论
//文章value:用户:文章
func UserStatKey(userid uint, article bool) string {
	if article {
		return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":article:stat:set")
	} else {

		return tool.StrSplicing("articlemanager:userid:", UintToStr(userid), ":comment:stat:set")
	}
}
func UserStatCommentValue(userid, artid, commentid uint) string {
	return tool.StrSplicing(UintToStr(userid), ":", UintToStr(artid), ":", UintToStr(commentid))
}
func UserStatArticleValue(userid, artid uint) string {
	return tool.StrSplicing(UintToStr(userid), ":", UintToStr(artid))
}
