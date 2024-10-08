package model

import "time"

type Tag struct {
	Model
	Name string `gorm:"unique;type:varchar(20);not null" json:"name"`

	// 如果该字段的值为零值（例如，对于切片或映射，它是nil；对于数字，它是0等），那么在序列化为JSON时这个字段会被忽略
	// joinForeignKey: 省略了，gorm默认外键是tag_id
	Articles []*Article `gorm:"many2many:article_tag;" json:"articles,omitempty"`
}
type TagVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}
