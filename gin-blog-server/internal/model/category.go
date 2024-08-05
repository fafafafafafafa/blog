package model

// hasMany: 一个分类下可以有多篇文章
type Category struct {
	Model
	Name     string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryId"`
}
