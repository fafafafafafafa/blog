package model

// hasMany: 一个分类下可以有多篇文章
type Category struct {
	Model
	Name     string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryId"`
}

type CategoryVO struct {
	Category
	ArticleCount int `json:"article_count"`
}

// 添加/编辑分类对象
type AddOrEditCategoryReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}
