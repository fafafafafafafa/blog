package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetRoleList(db *gorm.DB, num, size int, keyword string) (list []model.RoleVO, total int64, err error) {
	db = db.Model(&model.Role{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Select("id", "name", "label", "created_at", "is_disable").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func GetResourceIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&model.RoleResource{}).
		Where("role_id = ?", roleId).
		Pluck("resource_id", &ids)
	return ids, result.Error
}

func GetMenuIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&model.RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &ids)
	return ids, result.Error
}

func AddRole(db *gorm.DB, name, label string) error {
	role := model.Role{
		Name:  name,
		Label: label,
	}
	result := db.Create(&role)
	return result.Error
}

func UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	role := model.Role{
		Model:     model.Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(&role).Select("name", "label", "is_disable").Updates(&role).Error; err != nil {
			return err
		}

		// role_resource
		if err := db.Delete(&model.RoleResource{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		for _, rid := range resourceIds {
			if err := db.Create(&model.RoleResource{RoleId: role.ID, ResourceId: rid}).Error; err != nil {
				return err
			}
		}

		// role_menu
		if err := db.Delete(&model.RoleMenu{}, "role_id = ?", id).Error; err != nil {
			return err
		}

		for _, mid := range menuIds {
			if err := db.Create(&model.RoleMenu{RoleId: role.ID, MenuId: mid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 删除角色: 事务删除 role, role_resource, role_menu
func DeleteRoles(db *gorm.DB, ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {

		result := db.Delete(&model.Role{}, "id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&model.RoleResource{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&model.RoleMenu{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func GetRoleOption(db *gorm.DB) (list []model.OptionVO, err error) {
	result := db.Model(&model.Role{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}
