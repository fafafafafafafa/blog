package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetResource(db *gorm.DB, uri, method string) (resource model.Resource, err error) {
	result := db.Where(&model.Resource{Url: uri, Method: method}).First(&resource)
	return resource, result.Error
}

func CheckRoleAuth(db *gorm.DB, rid int, uri, method string) (bool, error) {
	resources, err := GetResourcesByRole(db, rid)
	if err != nil {
		return false, err
	}

	for _, r := range resources {
		if r.Anonymous || (r.Url == uri && r.Method == method) {
			return true, nil
		}
	}

	return false, nil
}

func GetResourcesByRole(db *gorm.DB, rid int) (resources []model.Resource, err error) {
	var role model.Role
	result := db.Model(&model.Role{}).Preload("Resources").Take(&role, rid)
	return role.Resources, result.Error
}
