package dao

import (
	"gin-blog/internal/model"
	"time"

	"gorm.io/gorm"
)

func GetUserAuthInfoByName(db *gorm.DB, name string) (*model.UserAuth, error) {
	var userAuth model.UserAuth
	result := db.Where(&model.UserAuth{Username: name}).First(&userAuth)
	return &userAuth, result.Error
}

func GetUserAuthInfoById(db *gorm.DB, id int) (*model.UserAuth, error) {
	var userAuth = model.UserAuth{Model: model.Model{ID: id}}
	result := db.Model(&userAuth).
		Preload("Roles").Preload("UserInfo").
		First(&userAuth)
	return &userAuth, result.Error
}

// 更新用户登录信息
func UpdateUserAuthLoginInfo(db *gorm.DB, id int, ipAddress, ipSource string) error {
	now := time.Now()
	userAuth := model.UserAuth{
		IpAddress:     ipAddress,
		IpSource:      ipSource,
		LastLoginTime: &now,
	}

	result := db.Where("id=?", id).Updates(userAuth)
	return result.Error
}

func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	result := db.
		Model(&model.UserAuthRole{UserAuthId: userAuthId}).
		Pluck("role_id", &ids)
	return ids, result.Error
}
