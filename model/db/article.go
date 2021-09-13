package db

import (
	. "blog/model"
	"blog/tool"
)

func ExistArticle(artid uint) bool {
	count := 0

	DB.Model(&Article{}).Where("id = ?", artid).Count(&count)
	if count == 0 {
		return false
	}
	return true
}
func ExistComment(commentid uint) bool {
	count := 0

	DB.Model(&Comment{}).Where("id = ?", commentid).Count(&count)
	if count == 0 {
		return false
	}
	return true

}
func CommentID(artid uint) ([]int64, error) {
	var commentids []int64
	if err := DB.Model(&Comment{}).Where("article_id=?", artid).Pluck("id", &commentids).Error; err != nil {
		return nil, err
	}
	return commentids, nil
}
func CommentList(commentids []int64) ([]Comment, error) {
	comments := make([]Comment, len(commentids))
	if err := DB.Where("id IN (?)", commentids).Find(&comments).Error; err != nil {

		return nil, err
	}
	for k, comment := range comments {

		var username []string
		if err := DB.Model(&User{}).Where("id=?", comment.UserID).Pluck("user_name", &username).Error; err != nil {
			return nil, err
		}
		comments[k].UserName = username[0]
		var count int64
		if err := DB.Model(&Stat{}).Where("type=? AND stat_id=? AND state=?", 1, comment.ID, 0).Count(&count).Error; err != nil {

			return nil, err
		}
		comments[k].Stat = uint(count)
	}
	return comments, nil
}
func DeleteArticle(articleid uint) error {
	//删除文章 1.删点赞 2.删评论 3.删文章
	//var (
	//stat    Stat
	//comment Comment
	//)
	//这里应该用事务
	tx := DB.Begin()
	if err := tx.Where("type=? AND stat_id=?", 0, articleid).Delete(&Stat{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("article_id=?", articleid).Delete(&Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&Article{}, articleid).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func DeleteComment(commentids []int64) error {
	//开启mysql事务
	tx := DB.Begin()
	//删评论及子评论点赞
	if err := tx.Where("type=? AND stat_id IN (?)", 1, commentids).Delete(&Stat{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//删评论及子评论
	//有问题
	if err := tx.Where("id IN (?)", commentids).Delete(&Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//队列dfs
func SubCommentid(commentid uint) ([]int64, error) {
	commentids := []int64{int64(commentid)}
	var queue []int64 //切片模拟队列
	//rootid加入队列
	if err := DB.Model(&Comment{}).Where("f_comment_id=?", commentid).Pluck("id", &queue).Error; err != nil {
		return nil, err
	}
	for len(queue) != 0 {
		var temp []int64
		if err := DB.Model(&Comment{}).Where("f_comment_id=?", queue[0]).Pluck("id", &temp).Error; err != nil {
			return nil, err
		}
		for _, v := range temp {
			queue = tool.Push(queue, v)
		}
		var cid int64
		cid, queue = tool.Pop(queue)
		commentids = append(commentids, cid)
	}
	return commentids, nil
}
