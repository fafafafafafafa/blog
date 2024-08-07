package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetMenuList(db *gorm.DB, keyword string) (list []model.Menu, total int64, err error) {
	db = db.Model(&model.Menu{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	result := db.Count(&total).Find(&list)
	return list, total, result.Error
}

// 获取所有菜单列表（超级管理员用）
func GetAllMenuList(db *gorm.DB) (menu []model.Menu, err error) {
	result := db.Find(&menu)
	return menu, result.Error
}

// 根据 user_id 获取菜单列表
func GetMenuListByUserId(db *gorm.DB, id int) (menus []model.Menu, err error) {
	var userAuth model.UserAuth
	result := db.Where(&model.UserAuth{Model: model.Model{ID: id}}).
		Preload("Roles").Preload("Roles.Menus").
		First(&userAuth)

	if result.Error != nil {
		return nil, result.Error
	}

	set := make(map[int]model.Menu)
	for _, role := range userAuth.Roles {
		for _, menu := range role.Menus {
			set[menu.ID] = menu
		}
	}

	for _, menu := range set {
		menus = append(menus, menu)
	}

	return menus, nil
}

func AddOrUpdateMenu(db *gorm.DB, menu *model.Menu) error {
	var result *gorm.DB

	if menu.ID > 0 {
		result = db.Model(menu).
			Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
			Updates(menu)
	} else {
		result = db.Create(menu)
	}

	return result.Error
}

func CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&model.RoleMenu{}).Where("menu_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func GetMenuById(db *gorm.DB, id int) (menu *model.Menu, err error) {
	result := db.First(&menu, id)
	return menu, result.Error
}

func CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func DeleteMenu(db *gorm.DB, id int) error {
	result := db.Delete(&model.Menu{}, id)
	return result.Error
}
