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

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set

)
