package redis

import (
	. "blog/model"
	"blog/tool"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

func AddArticle(uid uint, title, content string, tags []string) (int, error) {
	if err := Redisdb.SetNX(GetArticleIDKey(uid), 1, 0).Err(); err != nil && err != RedisNil {
		return -1, err
	}
	artid, err := Redisdb.Get(GetArticleIDKey(uid)).Int() //如果在事务外获取文章ID，会闭包逃逸
	if err != nil && err != RedisNil {
		return -1, err
	}
	transactional := func(tx *redis.Tx) error {
		article := map[string]interface{}{
			"title":   title,
			"content": content,
			"time":    tool.ShortTime(),
		}
		if tags != nil {
			article["tags"] = tool.SliceToString(tags)
		}

		//写入文章
		if err := tx.HMSet(ArticleIdKey(uid, uint(artid)), article).Err(); err != nil {
			return err
		}
		//写入用户文章集合
		if err := tx.SAdd(AuthorArticlesKey(uid), artid).Err(); err != nil {
			return err
		}
		//若设置了标签，使用集合存文章tag，和tag对应的文章
		//可以根据标签找到相应文章
		//写入文章tag集合
		if tags != nil {
			if err := tx.SAdd(ArticleTagsKey(uid, uint(artid)), tags).Err(); err != nil {
				return err
			}
			for _, tag := range tags {
				if tag != "" {
					//写入tag对应的文章
					if err := tx.SAdd(TagKey(uid, tag), artid).Err(); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
	if err := Redisdb.Watch(transactional, GetArticleIDKey(uid)); err != nil { //监控该用户文章ID主键,保证并发性
		return -1, err
	}
	if err := Redisdb.Incr(GetArticleIDKey(uid)).Err(); err != nil { //若事务提交成功了，自增该用户文章ID
		return -1, err
	}
	return artid, nil
}
func AddComment(who, uid, artid uint, content string) (int, error) {
	if err := Redisdb.SetNX(GetArticleCommentsIDKey(uid, artid), 1, 0).Err(); err != nil && err != RedisNil {
		return -1, err
	}
	commentid, err := Redisdb.Get(GetArticleCommentsIDKey(uid, artid)).Int()
	if err != nil && err != RedisNil {
		return -1, err
	}
	transactional := func(tx *redis.Tx) error {
		//写入评论集
		if err = tx.ZAdd(ArticleCommentRankKey(uid, artid), InitCommentRank(commentid)).Err(); err != nil {

			return err
		}
		comment := Comment{
			CommentId: uint(commentid),
			UserId:    who,
			AuthorId:  uid,
			Time:      tool.ShortTime(),
			Content:   content,
		}
		comment_json, err := json.Marshal(comment) //序列化
		if err != nil {
			return err
		}
		//写入评论
		if err = tx.HSet(ArticleCommentIDKey(uid, artid, uint(commentid)), "comment:0", comment_json).Err(); err != nil {

			return err
		}
		data, err := tx.HGet(ArticleCommentIDKey(uid, artid, uint(commentid)), "comment:0").Result()
		fmt.Println(data)
		fmt.Println(ArticleCommentIDKey(uid, artid, uint(commentid)))
		return nil
	}
	if err := Redisdb.Watch(transactional, GetArticleCommentsIDKey(uid, artid)); err != nil {
		return -1, err
	}
	if err := Redisdb.Incr(GetArticleCommentsIDKey(uid, artid)).Err(); err != nil {
		return -1, err
	}
	return commentid, nil
}
func StatComment(who, uid, artid, commentid uint) error {
	//点赞集合,防止重复点赞
	pipe := Redisdb.Pipeline()
	if ok, err := pipe.SAdd(UserStatKey(who), UserValue(uid, artid, commentid)).Result(); ok == 0 {
		return errors.New("aleady stat")
	} else if err != nil {
		return err
	}
	//点赞数+1
	if err := pipe.Incr(ArticleCommentStatKey(uid, artid, commentid)).Err(); err != nil {
		return err
	}
	//更新Rank
	if err := pipe.ZIncrBy(ArticleCommentRankKey(uid, artid), 1, strconv.Itoa(int(commentid))).Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}
