package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetPageList(db *gorm.DB) ([]model.Page, int64, error) {
	var pages = make([]model.Page, 0)
	var total int64

	result := db.Model(&model.Page{}).Count(&total).Find(&pages)
	return pages, total, result.Error
}

func AddOrUpdatePage(db *gorm.DB, id int, name, label, cover string) (*model.Page, error) {
	page := model.Page{
		Model: model.Model{ID: id},
		Name:  name,
		Label: label,
		Cover: cover,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&page)
	} else {
		result = db.Create(&page)
	}

	return &page, result.Error
}
