package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetResourceList(db *gorm.DB, keyword string) (list []model.Resource, err error) {
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	result := db.Find(&list)
	return list, result.Error
}

// Resource

func AddOrUpdateResource(db *gorm.DB, id, pid int, name, url, method string) error {
	resource := model.Resource{
		Model:    model.Model{ID: id},
		Name:     name,
		Url:      url,
		Method:   method,
		ParentId: pid,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&resource)
	} else {
		result = db.Create(&resource)
		// TODO: ????
		// * 解决前端的 BUG: 级联选中某个父节点后, 新增的子节点默认会展示被选中, 实际上未被选中值
		// * 解决方案: 新增子节点后, 删除该节点对应的父节点与角色的关联关系
		// dao.Delete(model.RoleResource{}, "resource_id", data.ParentId)
	}
	return result.Error
}

func DeleteResource(db *gorm.DB, id int) (int, error) {
	result := db.Delete(&model.Resource{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func CheckResourceInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&model.RoleResource{}).Where("resource_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func GetResourceById(db *gorm.DB, id int) (resource model.Resource, err error) {
	result := db.First(&resource, id)
	return resource, result.Error
}

func CheckResourceHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&model.Resource{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func UpdateResourceAnonymous(db *gorm.DB, id int, anonymous bool) error {
	result := db.Model(&model.Resource{}).Where("id = ?", id).Update("anonymous", anonymous)
	return result.Error
}
