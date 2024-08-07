package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

// 获取后台评论列表
func GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) (data []model.Comment, total int64, err error) {
	if typ != 0 {
		db = db.Where("type = ?", typ)
	}
	if isReview != nil {
		db = db.Where("is_review = ?", *isReview)
	}
	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	result := db.Model(&model.Comment{}).
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Preload("Article").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)

	return data, total, result.Error
}
