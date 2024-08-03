package models

import "github.com/jinzhu/gorm"

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	if err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return auth.ID > 0, nil
}
