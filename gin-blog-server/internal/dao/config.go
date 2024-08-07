package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetConfigMap(db *gorm.DB) (map[string]string, error) {
	var configs []model.Config
	result := db.Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	m := make(map[string]string)
	for _, config := range configs {
		m[config.Key] = config.Value
	}

	return m, nil
}

func CheckConfigMap(db *gorm.DB, m map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for k, v := range m {
			result := tx.Model(model.Config{}).Where("key=?", k).Update("value", v)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func GetConfig(db *gorm.DB, key string) string {
	var config model.Config
	result := db.Where("key", key).First(&config)

	if result.Error != nil {
		return ""
	}

	return config.Value
}

func CheckConfig(db *gorm.DB, key, value string) error {
	var config model.Config

	result := db.Where("key", key).FirstOrCreate(&config)
	if result.Error != nil {
		return result.Error
	}

	config.Value = value
	result = db.Save(&config)

	return result.Error
}
