package model

type Page struct {
	Model
	Name  string `gorm:"unique;type:varchar(20)" json:"name"`
	Label string `gorm:"unique;type:varchar(30)" json:"label"`
	Cover string `gorm:"type:varchar(255)" json:"cover"`
}
