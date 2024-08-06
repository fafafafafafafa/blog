package g

const (
	CTX_DB        = "_db_field"
	CTX_RDB       = "_rdb_field"
	CTX_USER_AUTH = "_user_auth_field"
)

// Redis Key
const (
	ONLINE_USER  = "online_user:"  // 在线用户
	OFFLINE_USER = "offline_user:" // 强制下线用户
	VISITOR_AREA = "visitor_area"  // 地域统计
	VIEW_COUNT   = "view_count"    // 访问数量

	KEY_UNIQUE_VISITOR_SET = "unique_visitor" // 唯一用户记录 set

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set

	CONFIG = "config" // 博客配置

)
