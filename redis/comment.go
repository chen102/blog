package redis

import (
	"blog/model"
	"blog/tool"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func ShowCommentListCache(artid, offset, count uint) ([]model.Comment, error) {
	ok, err := model.RedisReadDB.Exists(ArticleCommentKey(artid)).Result()
	if err != nil && err != model.RedisNil {
		return nil, err
	} else if ok == 0 {
		return nil, model.RedisNil
	}
	ids, err := model.RedisReadDB.SMembers(ArticleCommentKey(artid)).Result()
	if err != nil {
		return nil, err
	}

	get := GetSort(CommentStringKey("*"), true, "UserID", "UserName", "RevUserName", "UpdatedAt", "Stat", "Content", "FCommentID") //排序后显示的字段

	sortargs := SortArgs(tool.StrSplicing(CommentStringKey("*"), "->Stat"), int64(offset), int64(count), get, "DESC", false) //根据点赞数排序
	data, err := model.RedisWriteDB.Sort(ArticleCommentKey(artid), sortargs).Result()
	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, errors.New("没拿到排序结果")
	}
	comments, err := model.CommentRank(data)
	if err != nil {
		return nil, err
	}
	//刷新缓存时间
	for _, id := range ids {
		if err := model.RedisWriteDB.Expire(CommentStringKey(id), 1*time.Hour).Err(); err != nil {
			return nil, err
		}
	}
	//page := tool.PaginationINT(ids, int(offset), int(count))
	//if page == nil {
	//return nil, nil
	//}
	//comments := make([]model.Comment, len(page))
	//for k, commentid := range page {
	//data, err := model.RedisReadDB.HGetAll(CommentKey(uint(commentid))).Result()
	//if err != nil {
	//return nil, err
	//} else if data != nil {
	//mapstructure.WeakDecode(data, &comments[k])
	//}
	//model.RedisWriteDB.Expire(CommentKey(uint(commentid)), 1*time.Hour)
	return comments, nil
}

func WriteCommentListCache(artid uint, comments []model.Comment) error {
	transactional := func(tx *redis.Tx) error {
		for _, comment := range comments {
			if err := tx.SAdd(ArticleCommentKey(artid), comment.ID).Err(); err != nil {
				return nil
			}
			if err := tx.HMSet(CommentKey(comment.ID), model.StructToMap(comment)).Err(); err != nil {
				return err
			}
			tx.Expire(CommentKey(comment.ID), 1*time.Hour)

		}
		tx.Expire(ArticleCommentKey(artid), 1*time.Hour)

		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticleCommentKey(artid)); err != nil {
		return err
	}
	return nil

}
func DeleteComment(artid uint) error {
	if err := model.RedisWriteDB.Del(ArticleCommentKey(artid)).Err(); err != nil {
		return err
	}
	return nil
}
func CommentIncrementCache(artid uint, commentids []int64) ([]int64, error) {
	transactional := func(tx *redis.Tx) error {
		for k, commentid := range commentids {
			if err := tx.SAdd(ArticleCommentKey(artid), commentid).Err(); err != nil {
				return err
			}
			exist, err := tx.Exists(CommentKey(uint(commentid))).Result()
			if err != nil && err != model.RedisNil {
				return err
			} else if exist == 1 {
				if err := tx.Expire(CommentKey(uint(commentid)), 1*time.Hour).Err(); err != nil {
					return err
				}
				commentids[k] = -1
			}

		}
		return nil
	}
	if err := model.RedisWriteDB.Watch(transactional, ArticleCommentKey(artid)); err != nil {
		return nil, err
	}
	return commentids, nil
}
