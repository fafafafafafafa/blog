package model

// belongTo: 一个文章 属于 一个分类
// belongTo: 一个文章 属于 一个用户
// many2many: 一个文章 可以拥有 多个标签, 多个文章 可以使用 一个标签
type Article struct {
	Model

	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"` // 1-原创 2-转载 3-翻译
	Status      int    `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`    // 1-公开 2-私密
	IsTop       bool   `json:"is_top"`
	IsDelete    bool   `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	// article_tag 这是用于连接两个表的中间表的名称。在 GORM 中处理多对多关系时，通常需要一个中间表来存储两个表之间的关联信息。这个中间表通常包含两个外键，分别指向两个相关联的表
	// article_id: 这部分指定了连接外键为 article_id，也就是说在 article_tag 这个中间表中，用于关联文章的字段名称是 article_id
	Tags []*Tag `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`

	CategoryId int       `json:"category_id"`
	Category   *Category `gorm:"foreignkey:CategoryId" json:"category"`

	UserId int       `json:"-"` // user_auth_id
	User   *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

// TODO: 添加对标签数组的查询
type ArticleQuery struct {
	PageQuery
	Title      string `form:"title"`
	CategoryId int    `form:"category_id"`
	TagId      int    `form:"tag_id"`
	Type       int    `form:"type"`
	Status     int    `form:"status"`
	IsDelete   *bool  `form:"is_delete"`
}

type ArticleVO struct {
	Article

	LikeCount    int `json:"like_count" gorm:"-"`
	ViewCount    int `json:"view_count" gorm:"-"`
	CommentCount int `json:"comment_count" gorm:"-"`
}

type AddOrEditArticleReq struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Desc        string `json:"desc"`
	Content     string `json:"content" binding:"required"`
	Img         string `json:"img"`
	Type        int    `json:"type" binding:"required,min=1,max=3"`   // 类型: 1-原创 2-转载 3-翻译
	Status      int    `json:"status" binding:"required,min=1,max=3"` // 状态: 1-公开 2-私密 3-评论可见
	IsTop       bool   `json:"is_top"`
	OriginalUrl string `json:"original_url"`

	TagNames     []string `json:"tag_names"`
	CategoryName string   `json:"category_name"`
}

type UpdateArticleTopReq struct {
	ID    int  `json:"id"`
	IsTop bool `json:"is_top"`
}

type SoftDeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
}
