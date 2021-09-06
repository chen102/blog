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
func CommentList(artid uint) ([]Comment, error) {
	var comments []Comment
	if err := DB.Where("article_id=?", artid).Find(&comments).Error; err != nil {

		return nil, err
	}
	for k, comment := range comments {

		var username []string
		if err := DB.Model(&User{}).Where("id=?", comment.UserID).Pluck("user_name", &username).Error; err != nil {
			return nil, err
		}
		comments[k].UserName = username[0]
		var count int64
		if err := DB.Model(&Stat{}).Where("type=? AND article_id=? AND Stat=?", 1, comment.ID, 0).Count(&count).Error; err != nil {

			return nil, err
		}
		comments[k].Stat = uint(count)
	}
	return comments, nil
}
