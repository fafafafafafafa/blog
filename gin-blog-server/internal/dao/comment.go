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

// 获取博客评论列表
func GetCommentVOList(db *gorm.DB, page, size, topic, typ int) (data []model.CommentVO, total int64, err error) {
	var list []model.Comment

	tx := db.Model(&model.Comment{})
	if typ != 0 {
		tx = tx.Where("type = ?", typ)
	}
	if topic != 0 {
		tx = tx.Where("topic_id = ?", topic)
	}

	// 获取顶级评论列表
	tx.Where("parent_id = 0").
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		// Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size))
	if err := tx.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 获取顶级评论的回复列表
	for _, v := range list {
		replyList := make([]model.CommentVO, 0)

		tx := db.Model(&model.Comment{})
		tx.Where("parent_id = ?", v.ID).
			Preload("User").Preload("User.UserInfo").
			// Preload("ReplyUser").Preload("ReplyUser.UserInfo")
			Order("id DESC")
		if err := tx.Find(&replyList).Error; err != nil {
			return nil, 0, err
		}

		data = append(data, model.CommentVO{
			ReplyCount: len(replyList),
			Comment:    v,
			ReplyList:  replyList,
		})
	}

	return data, total, nil
}

// 根据 [评论id] 获取 [回复列表]
func GetCommentReplyList(db *gorm.DB, id, page, size int) (data []model.Comment, err error) {
	result := db.Model(&model.Comment{}).
		Where(&model.Comment{ParentId: id}).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)
	return data, result.Error
}

// 新增评论
func AddComment(db *gorm.DB, userId, typ, topicId int, content string, isReview bool) (*model.Comment, error) {
	comment := model.Comment{
		UserId:   userId,
		TopicId:  topicId,
		Content:  content,
		Type:     typ,
		IsReview: isReview,
	}
	result := db.Create(&comment)
	return &comment, result.Error
}

// 回复评论
func ReplyComment(db *gorm.DB, userId, replyUserId, parentId int, content string, isReview bool) (*model.Comment, error) {
	var parent model.Comment
	result := db.First(&parent, parentId)
	if result.Error != nil {
		return nil, result.Error
	}

	comment := model.Comment{
		UserId:      userId,
		Content:     content,
		ReplyUserId: replyUserId,
		ParentId:    parentId,
		IsReview:    isReview,
		TopicId:     parent.TopicId, // 主题和父评论一样
		Type:        parent.Type,    // 类型和父评论一样
	}
	result = db.Create(&comment)
	return &comment, result.Error
}
