package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func UpdateMessagesReview(db *gorm.DB, ids []int, isReview bool) (int64, error) {
	result := db.Model(&model.Message{}).Where("id in ?", ids).Update("is_review", isReview)
	return result.RowsAffected, result.Error
}
func GetMessageList(db *gorm.DB, num, size int, nickname string, isReview *bool) (list []model.Message, total int64, err error) {
	db = db.Model(&model.Message{})

	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	if isReview != nil {
		db = db.Where("is_review = ?", isReview)
	}

	db.Count(&total)
	result := db.Order("created_at DESC").Scopes(Paginate(num, size)).Find(&list)
	return list, total, result.Error
}

func DeleteMessages(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id in ?", ids).Delete(&model.Message{})
	return result.RowsAffected, result.Error
}

func AddMessage(db *gorm.DB, nickname, avatar, content, address, source string, speed int, isReview bool) (*model.Message, error) {
	message := model.Message{
		Nickname:  nickname,
		Avatar:    avatar,
		Content:   content,
		IpAddress: address,
		IpSource:  source,
		Speed:     speed,
		IsReview:  isReview,
	}

	result := db.Create(&message)
	return &message, result.Error
}
