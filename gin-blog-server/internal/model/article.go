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
