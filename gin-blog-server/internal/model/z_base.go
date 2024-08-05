package model

import (
	"time"

	"gorm.io/gorm"
)

// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) error {
	// 设置表关联
	db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})
	db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})
	//! article 和 tag 不需要设置表关联吗
	// db.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})

	return db.AutoMigrate(
		&Article{},      // 文章
		&Category{},     // 分类
		&Tag{},          // 标签
		&Comment{},      // 评论
		&Message{},      // 消息
		&FriendLink{},   // 友链
		&Page{},         // 页面
		&Config{},       // 网站设置
		&OperationLog{}, // 操作日志
		&UserInfo{},     // 用户信息

		&UserAuth{},     // 用户验证
		&Role{},         // 角色
		&Menu{},         // 菜单
		&Resource{},     // 资源（接口）
		&RoleMenu{},     // 角色-菜单 关联
		&RoleResource{}, // 角色-资源 关联
		&UserAuthRole{}, // 用户-角色 关联
	)
}

// 通用模型

type Model struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
