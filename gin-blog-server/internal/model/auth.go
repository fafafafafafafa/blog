package model

import (
	"time"
)

// 权限控制: 7 张表（4 模型 + 3 关联）
// belong to: 一个UserAuth用户帐号 属于一个 UserInfo用户信息
// many2many: 一个UserAuth用户帐号 拥有多个类型的Role角色
type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`

	Roles []*Role `json:"roles" gorm:"many2many:user_auth_role"`
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

// many2many: 一个类型的Role角色 拥有多个 resource资源模块
// many2many: 一个类型的Role角色 拥有多个 menus菜单模块
// many2many: 一个类型的Role角色 拥有多个 UserAuth用户帐号
type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

// many2many: 一个resource资源模块 拥有多个类型的Role角色
type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"`
}

/*
菜单设计:

目录: catalogue === true
  - 如果是目录, 作为单独项, 不展开子菜单（例如 "首页", "个人中心"）
  - 如果不是目录, 且 parent_id 为 0, 则为一级菜单, 可展开子菜单（例如 "文章管理" 下有 "文章列表", "文章分类", "文章标签" 等子菜单）
  - 如果不是目录, 且 parent_id 不为 0, 则为二级菜单

隐藏: hidden
  - 隐藏则不显示在菜单栏中

外链: external, external_link
  - 如果是外链, 如果设置为外链, 则点击后会在新窗口打开
*/
// many2many: 一个menus菜单模块 拥有多个类型的Role角色
type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"` // 菜单名称
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"` // 路由地址
	Component    string `gorm:"type:varchar(50)" json:"component"`                          // 组件路径
	Icon         string `gorm:"type:varchar(50)" json:"icon"`                               // 图标
	OrderNum     int8   `json:"order_num"`                                                  // 排序
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                           // 重定向地址
	Catalogue    bool   `json:"is_catalogue"`                                               // 是否为目录
	Hidden       bool   `json:"is_hidden"`                                                  // 是否隐藏
	KeepAlive    bool   `json:"keep_alive"`                                                 // 是否缓存
	External     bool   `json:"is_external"`                                                // 是否外链
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                     // 外链地址

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}
