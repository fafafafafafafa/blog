package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetUserInfoById(db *gorm.DB, id int) (*model.UserInfo, error) {
	var userInfo model.UserInfo
	result := db.Model(&userInfo).Where("id=?", id).First(&userInfo)
	return &userInfo, result.Error
}
