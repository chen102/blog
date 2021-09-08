//Redis所有Key
//命名规范:   模块名:业务逻辑:存储的东西:value类型
package redis

import (
	"blog/tool"
	"github.com/go-redis/redis"
	"strconv"
)

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

//文章键 articleservice:articleid:xx:article:hash
//value:用户的文章内容
func ArticleIdKey(articleid uint) string {

	return tool.StrSplicing("articleservice:userid:articleid:", UintToStr(articleid), ":article:hash")
}

func ArticleStringIdKey(articleid string) string {

	return tool.StrSplicing("articleservice:userid:articleid:", articleid, ":article:hash")
}

//用户键
func UserIdKey(userid uint) string {
	return tool.StrSplicing("userservice:userid:", UintToStr(userid), ":user:hash")
}
func UserStringIdKey(userid string) string {
	return tool.StrSplicing("userservice:userid:", userid, ":user:hash")
}

//文章列表键
func ArticlesListKey(userid uint) string {
	return tool.StrSplicing("articleservice:userid:", UintToStr(userid), ":articleid:set")
}

//用户点赞队列键
func UserStatQueueKey() string {
	return "statservice:redistomysql:userid:articleid:list"
}

//点赞队列消息
//0:点赞文章
//1:取消点赞文章
//2:点赞评论
//3.取消点赞评论
//用户点赞文章消息
func UserStatQueueValue(userid, artid uint) string {
	return tool.StrSplicing("0:", UintToStr(userid), ":", UintToStr(artid))
}

//用户取消点赞文章消息
func UserCancesStatQueueValue(userid, artid uint) string {
	return tool.StrSplicing("1:", UintToStr(userid), ":", UintToStr(artid))
}

//用户点赞评论消息
func StatCommenQueueValue(userid, commentid uint) string {
	return tool.StrSplicing("2:", UintToStr(userid), ":", UintToStr(commentid))
}

//用户取消点赞评论消息
func CancesStatCommenQueueValue(userid, commentid uint) string {
	return tool.StrSplicing("3:", UintToStr(userid), ":", UintToStr(commentid))
}

//用户点赞列表键
func UserStatList(userid uint) string {
	return tool.StrSplicing("statservice:userstatlist:userid:", UintToStr(userid), ":artid:set")
}

//用户关注列表键
func UserFolloweListrKey(userid, listtype uint) string {
	if listtype == 0 {
		return tool.StrSplicing("followerservice:followerlist:userid:", UintToStr(userid), ":userid:set")
	} else if listtype == 1 {
		return tool.StrSplicing("followerservice:fanslist:userid:", UintToStr(userid), ":userid:set")
	} else if listtype == 2 {
		return tool.StrSplicing("followerservice:mutualfollowing:", UintToStr(userid), ":userid:set")
	}
	return ""
}

//用户动态键
func UserDynamicKey(userid uint) string {
	return tool.StrSplicing("followerservice:dynamic:userid:", UintToStr(userid), ":articlid:set")
}

//文章评论键
func ArticleCommentKey(artid uint) string {
	return tool.StrSplicing("articleservice:comment:articleid:", UintToStr(artid), ":commentid:set")
}

//评论键
func CommentKey(commentid uint) string {
	return tool.StrSplicing("commentservice:commentid:", UintToStr(commentid), ":comment:hash")
}

//评论键
func CommentStringKey(commentid string) string {
	return tool.StrSplicing("commentservice:commentid:", commentid, ":comment:hash")
}

//----------------------------------------------
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
