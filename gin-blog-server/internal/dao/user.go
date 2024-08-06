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

func GetUserList(db *gorm.DB, page, size int, loginType int8, nickname, username string) (list []model.UserAuth, total int64, err error) {
	if loginType != 0 {
		db = db.Where("login_type = ?", loginType)
	}

	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}

	result := db.Model(&model.UserAuth{}).
		Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
		Where("user_info.nickname LIKE ?", "%"+nickname+"%").
		Preload("UserInfo").
		Preload("Roles").
		Count(&total).
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

// 更新用户昵称及角色信息
func UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error {
	// 开启事物
	return db.Transaction(func(tx *gorm.DB) error {
		userAuth, err := GetUserAuthInfoById(db, authId)
		if err != nil {
			return err
		}

		userInfo := model.UserInfo{
			Model:    model.Model{ID: userAuth.UserInfoId},
			Nickname: nickname,
		}
		result := db.Model(&userInfo).Updates(userInfo)
		if result.Error != nil {
			return result.Error
		}

		// 至少有一个角色
		if len(roleIds) == 0 {
			return nil
		}

		// 更新用户角色, 清空原本的 user_role 关系, 添加新的关系
		result = db.Where(model.UserAuthRole{UserAuthId: userAuth.UserInfoId}).Delete(model.UserAuthRole{})
		if result.Error != nil {
			return result.Error
		}

		var userRoles []model.UserAuthRole
		for _, id := range roleIds {
			userRoles = append(userRoles, model.UserAuthRole{
				RoleId:     id,
				UserAuthId: userAuth.ID,
			})
		}
		result = db.Create(&userRoles)

		return result.Error
	})

}
func UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	userAuth := model.UserAuth{
		Model:     model.Model{ID: id},
		IsDisable: isDisable,
	}
	result := db.Model(&userAuth).Select("is_disable").Updates(&userAuth)
	return result.Error
}

func UpdateUserPassword(db *gorm.DB, id int, password string) error {
	userAuth := model.UserAuth{
		Model:    model.Model{ID: id},
		Password: password,
	}
	result := db.Model(&userAuth).Select("password").Updates(&userAuth)
	return result.Error
}

func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	userInfo := model.UserInfo{
		Model:    model.Model{ID: id},
		Nickname: nickname,
		Avatar:   avatar,
		Intro:    intro,
		Website:  website,
	}

	result := db.Model(&userInfo).
		Select("nickname", "avatar", "intro", "website").
		Updates(userInfo)
	return result.Error
}
