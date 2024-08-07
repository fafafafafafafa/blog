package handle

import (
	"gin-blog/internal/dao"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Front struct{}

// 前台首页信息
func (*Front) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	data, err := dao.GetFrontStatistics(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	data.ViewCount, _ = rdb.Get(rctx, g.VIEW_COUNT).Int64()

	ReturnSuccess(c, data)
}

/*
文章相关接口
*/

// 获取文章列表
func (*Front) GetArticleList(c *gin.Context) {
	var query model.FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, _, err := dao.GetBlogArticleList(GetDB(c), query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

// 根据 [文章id] 获取 [文章详情]
func (*Front) GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	// 文章详情
	val, err := dao.GetBlogArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	article := model.BlogArticleVO{Article: *val}

	// 推荐文章（6篇）
	article.RecommendArticles, err = dao.GetRecommendList(db, id, 6)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 最新文章（5篇）
	article.NewestArticles, err = dao.GetNewestList(db, 5)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 更新文章浏览量 TODO: 删除文章时删除其浏览量
	// updateArticleViewCount(c, id)

	// TODO: 更新访问量
	// * 目前请求一次就会增加访问量, 即刷新可以刷访问量
	rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 1, strconv.Itoa(id))

	// 上一篇文章
	article.LastArticle, err = dao.GetLastArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 下一篇文章
	article.NextArticle, err = dao.GetNextArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 点赞量, 浏览量
	article.ViewCount = int64(rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, strconv.Itoa(id)).Val())
	article.LikeCount = int64(rdb.ZScore(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())

	// 评论数量
	article.CommentCount, err = dao.GetArticleCommentCount(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// 获取文章归档
func (*Front) GetArchiveList(c *gin.Context) {
	var query model.FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := dao.GetBlogArticleList(GetDB(c), query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	archives := make([]model.ArchiveVO, 0)
	for _, article := range list {
		archives = append(archives, model.ArchiveVO{
			ID:         article.ID,
			Title:      article.Title,
			Created_at: article.CreatedAt,
		})
	}

	ReturnSuccess(c, model.PageResult[model.ArchiveVO]{
		Total: total,
		List:  archives,
		Page:  query.Page,
		Size:  query.Size,
	})
}

// 获取友链列表
func (*Front) GetLinkList(c *gin.Context) {
	list, _, err := dao.GetLinkList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

// 查询消息列表
func (*Front) GetMessageList(c *gin.Context) {
	isReview := true
	list, _, err := dao.GetMessageList(GetDB(c), 1, 1000, "", &isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// 查询分类列表
func (*Front) GetCategoryList(c *gin.Context) {
	list, _, err := dao.GetCategoryList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// 查询标签列表
func (*Front) GetTagList(c *gin.Context) {
	list, _, err := dao.GetTagList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// 文章搜索
func (*Front) SearchArticle(c *gin.Context) {
	result := make([]model.ArticleSearchVO, 0)

	keyword := c.Query("keyword")
	if keyword == "" {
		ReturnSuccess(c, result)
		return
	}

	db := GetDB(c)

	articleList, err := dao.List(db, []model.Article{}, "*", "",
		"is_delete = 0 AND status = 1 AND (title LIKE ? OR content LIKE ?)",
		"%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	for _, article := range articleList {
		// 高亮标题中的关键字
		title := strings.ReplaceAll(article.Title, keyword,
			"<span style='color:#f47466'>"+keyword+"</span>")

		content := article.Content
		// 关键字在内容中的起始位置
		keywordStartIndex := unicodeIndex(content, keyword)
		if keywordStartIndex != -1 { // 关键字在内容中
			preIndex, afterIndex := 0, 0
			if keywordStartIndex > 25 {
				preIndex = keywordStartIndex - 25
			}
			// 防止中文截取出乱码 (中文在 golang 是 3 个字符, 使用 rune 中文占一个数组下标)
			preText := substring(content, preIndex, keywordStartIndex) // 截取关键词最多前25的unicode字符
			// string([]rune(content[preIndex:keywordStartIndex]))

			// 关键字在内容中的结束位置
			keywordEndIndex := keywordStartIndex + unicodeLen(keyword)
			// 截取关键词后最多175个unicode字符
			afterLength := len(content) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			// afterText := string([]rune(content)[keywordStartIndex:afterIndex])
			afterText := substring(content, keywordStartIndex, afterIndex)
			// 高亮内容中的关键字
			content = strings.ReplaceAll(preText+afterText, keyword,
				"<span style='color:#f47466'>"+keyword+"</span>")
		}

		result = append(result, model.ArticleSearchVO{
			ID:      article.ID,
			Title:   title,
			Content: content,
		})
	}

	ReturnSuccess(c, result)
}

// 获取带中文的字符串实际长度，非字节长度
func unicodeLen(str string) int {
	var r = []rune(str)
	return len(r)
}

// 获取带中文的字符串中子字符串的实际位置，非字节位置
func unicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

// 解决中文获取位置不正确问题
func substring(source string, start int, end int) string {
	var unicodeStr = []rune(source)
	length := len(unicodeStr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	// var substring = ""
	// for i := start; i < end; i++ {
	// 	substring += string(unicodeStr[i])
	// }
	// return substring
	return string(unicodeStr[start:end])
}

// 获取评论列表
func (*Front) GetCommentList(c *gin.Context) {
	var query model.FCommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	data, total, err := dao.GetCommentVOList(db, query.Page, query.Size, query.TopicId, query.Type)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	likeCountMap := rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val()
	for i, comment := range data {
		if len(data[i].ReplyList) > 3 {
			data[i].ReplyList = data[i].ReplyList[:3] // 只显示 3 条回复
		}
		data[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(comment.ID)])
	}

	ReturnSuccess(c, model.PageResult[model.CommentVO]{
		List:  data,
		Total: total,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// 根据 [评论id] 获取 [回复列表]
func (*Front) GetReplyListByCommentId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	var query model.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	replyList, err := dao.GetCommentReplyList(db, id, query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	likeCountMap := rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val()

	data := make([]model.CommentVO, 0)
	for _, reply := range replyList {
		like, _ := strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
		data = append(data, model.CommentVO{
			Comment:   reply,
			LikeCount: like,
		})
	}

	ReturnSuccess(c, data)
}

/*
以下接口需要登录
*/

// TODO: 添加自定义头像和昵称留言功能（即可以不登录留言）
// 保存留言（只能新增，不能编辑）
func (*Front) AddMessage(c *gin.Context) {
	var req model.FAddMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	db := GetDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)
	isReview := dao.GetConfigBool(db, g.CONFIG_IS_COMMENT_REVIEW)

	info := auth.UserInfo
	message, err := dao.AddMessage(db, info.Nickname, info.Avatar, req.Content, ipAddress, ipSource, req.Speed, isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, message)
}

// 保存评论（只能新增，不能编辑）
// TODO: 添加自定义头像和昵称留言功能（即可以不登录评论）
// TODO: 开启邮箱通知用户功能
// TODO: HTMLUtil.Filter 过滤 HTML 元素中的字符串...
func (*Front) AddComment(c *gin.Context) {
	var req model.FAddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	db := GetDB(c)
	isReview := dao.GetConfigBool(db, g.CONFIG_IS_COMMENT_REVIEW)

	var comment *model.Comment
	var err error

	if req.ReplyUserId == 0 {
		comment, err = dao.AddComment(db, auth.ID, req.Type, req.TopicId, req.Content, isReview)
	} else {
		comment, err = dao.ReplyComment(db, auth.ID, req.ReplyUserId, req.ParentId, req.Content, isReview)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, comment)
}

// 点赞评论
func (*Front) LikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)
	auth, _ := CurrentUserAuth(c)

	// 记录某个用户已经对某个评论点过赞
	commentLikeUserKey := g.COMMENT_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// 该评论已经被记录过, 再点赞就是取消点赞
	if rdb.SIsMember(rctx, commentLikeUserKey, id).Val() {
		rdb.SRem(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), -1)
	} else { // 未被记录过, 则是增加点赞
		rdb.SAdd(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), 1)
	}

	ReturnSuccess(c, nil)
}

// 点赞文章
// 需要记录某个用户已经对某篇文章点过赞, 防止重复点赞
func (*Front) LikeArticle(c *gin.Context) {
	auth, _ := CurrentUserAuth(c)

	articleId, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)

	// 记录某个用户已经对某个文章点过赞
	articleLikeUserKey := g.ARTICLE_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// 该文章已经被记录过, 再点赞就是取消点赞
	if rdb.SIsMember(rctx, articleLikeUserKey, articleId).Val() {
		rdb.SRem(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), -1)
	} else { // 未被记录过, 则是增加点赞
		rdb.SAdd(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), 1)
	}

	ReturnSuccess(c, nil)
}
