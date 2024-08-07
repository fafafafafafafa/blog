package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetOperationLogList(db *gorm.DB, num, size int, keyword string) (data []model.OperationLog, total int64, err error) {
	db = db.Model(&model.OperationLog{})
	if keyword != "" {
		db = db.Where("opt_module LIKE ?", "%"+keyword+"%").
			Or("opt_desc LIKE ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Order("created_at DESC").
		Scopes(Paginate(num, size)).
		Find(&data)
	return data, total, result.Error
}
