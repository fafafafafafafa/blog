package model

type FriendLink struct {
	Model
	Name    string `gorm:"type:varchar(50)" json:"name"`
	Avatar  string `gorm:"type:varchar(255)" json:"avatar"`
	Address string `gorm:"type:varchar(255)" json:"address"`
	Intro   string `gorm:"type:varchar(255)" json:"intro"`
}

// 添加或修改友链
type AddOrEditLinkReq struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Avatar  string `json:"avatar"`
	Address string `json:"address" binding:"required"`
	Intro   string `json:"intro"`
}
