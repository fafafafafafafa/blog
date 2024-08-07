package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetLinkList(db *gorm.DB, num, size int, keyword string) (list []model.FriendLink, total int64, err error) {
	db = db.Model(&model.FriendLink{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
		db = db.Or("address LIKE ?", "%"+keyword+"%")
		db = db.Or("intro LIKE ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Order("created_at DESC").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func AddOrUpdateLink(db *gorm.DB, id int, name, avatar, address, intro string) (*model.FriendLink, error) {
	link := model.FriendLink{
		Model:   model.Model{ID: id},
		Name:    name,
		Avatar:  avatar,
		Address: address,
		Intro:   intro,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&link)
	} else {
		result = db.Create(&link)
	}

	return &link, result.Error
}
