package model

const (
	TYPE_ARTICLE = iota + 1 // 文章
	TYPE_LINK               // 友链
	TYPE_TALK               // 说说
)

/*
如果评论类型是文章，那么 topic_id 就是文章的 id
如果评论类型是友链，不需要 topic_id
*/

type Comment struct {
	Model

	ParentId int    `json:"parent_id"` // 父评论
	Content  string `gorm:"type:varchar(500);not null" json:"content"`
	Type     int    `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
	IsReview bool   `json:"is_review"`

	// Belongs To
	UserId int       `json:"user_id"` // 评论者
	User   *UserAuth `gorm:"foreignKey:UserId" json:"user"`

	ReplyUserId int       `json:"reply_user_id"` // 被回复者
	ReplyUser   *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`

	TopicId int      `json:"topic_id"` // 评论的文章
	Article *Article `gorm:"foreignKey:TopicId" json:"article"`
}

type CommentQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
	Type     int    `form:"type"`
}
