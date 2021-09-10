package db

import . "blog/model"

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
func DeleteComment(commentid uint) error {
	tx := DB.Begin()
	//删评论点赞
	if err := tx.Where("type=? AND stat_id=?", 1, commentid).Delete(&Stat{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//删子评论
	if err := tx.Where("root_id=?", commentid).Delete(&Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//删评论
	if err := tx.Delete(&Comment{}, commentid).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
