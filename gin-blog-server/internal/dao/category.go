package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetCategoryList(db *gorm.DB, num, size int, keyword string) ([]model.CategoryVO, int64, error) {
	var list = make([]model.CategoryVO, 0)
	var total int64

	db = db.Table("category c").
		Select("c.id", "c.name", "COUNT(a.id) AS article_count", "c.created_at", "c.updated_at").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	result := db.Group("c.id").
		Order("c.updated_at DESC").
		Count(&total).
		Scopes(Paginate(num, size)).
		Find(&list)

	return list, total, result.Error
}

func AddOrUpdateCategory(db *gorm.DB, id int, name string) (*model.Category, error) {
	category := model.Category{
		Model: model.Model{ID: id},
		Name:  name,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(category)
	} else {
		result = db.Create(&category)
	}

	return &category, result.Error
}

func DeleteCategory(db *gorm.DB, ids []int) (int64, error) {
	// result := db.Where("id IN ?", ids).Delete(model.Category{})
	result := db.Model(&model.Category{}).Delete("id IN ?", ids)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func GetCategoryOption(db *gorm.DB) ([]model.OptionVO, error) {
	var list []model.OptionVO
	result := db.Model(&model.Category{}).Select("id", "name").Find(&list)
	return list, result.Error
}
