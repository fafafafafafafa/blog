package handle

import (
	"gin-blog/internal/dao"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"io"
	"log/slog"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Article struct{}

func (*Article) GetList(c *gin.Context) {
	var query model.ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	list, total, err := dao.GetArticleList(db, query.Page, query.Size, query.Title, query.IsDelete, query.Status, query.Type, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	likeCountMap := rdb.HGetAll(rctx, g.ARTICLE_LIKE_COUNT).Val()
	// 从一个有序集合中获取指定范围内的元素及其分数
	// 0 表示从第一个元素开始，-1 表示到最后一个元素（即获取整个集合的所有元素
	viewCountZ := rdb.ZRangeWithScores(rctx, g.ARTICLE_VIEW_COUNT, 0, -1).Val()

	viewCountMap := make(map[int]int)
	for _, article := range viewCountZ {
		id, _ := strconv.Atoi(article.Member.(string))
		viewCountMap[id] = int(article.Score)
	}

	data := make([]model.ArticleVO, 0)
	for _, article := range list {
		likeCount, _ := strconv.Atoi(likeCountMap[strconv.Itoa(article.ID)])
		data = append(data, model.ArticleVO{
			Article:   article,
			LikeCount: likeCount,
			ViewCount: viewCountMap[article.ID],
		})
	}

	ReturnSuccess(c, model.PageResult[model.ArticleVO]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  data,
	})

}

func (*Article) AddOrUpdate(c *gin.Context) {
	var req model.AddOrEditArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	if req.Img == "" {
		req.Img = dao.GetConfig(db, g.CONFIG_ARTICLE_COVER) // 默认图片
	}

	if req.Type == 0 {
		req.Type = 1 // 默认为原创
	}

	article := model.Article{
		Model:       model.Model{ID: req.ID},
		Title:       req.Title,
		Desc:        req.Desc,
		Content:     req.Content,
		Img:         req.Img,
		Type:        req.Type,
		Status:      req.Status,
		OriginalUrl: req.OriginalUrl,
		IsTop:       req.IsTop,
		UserId:      auth.UserInfoId,
	}

	err := dao.AddOrUpdateArticle(db, &article, req.CategoryName, req.TagNames)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// 修改置顶信息
func (*Article) UpdateTop(c *gin.Context) {
	var req model.UpdateArticleTopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := dao.UpdateArticleTop(GetDB(c), req.ID, req.IsTop)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 获取文章详细信息
func (*Article) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	article, err := dao.GetArticle(GetDB(c), id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

func (*Article) UpdateSoftDelete(c *gin.Context) {
	var req model.SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := dao.UpdateArticleSoftDelete(GetDB(c), req.Ids, req.IsDelete)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

func (*Article) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := dao.DeleteArticle(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// TODO: 目前是前端导出
// 导出文章: 获取导出后的资源链接列表
func (*Article) Export(c *gin.Context) {
	ReturnSuccess(c, nil)
}

// 导入文章: 题目 + 内容
func (*Article) Import(c *gin.Context) {
	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	fileName := fileHeader.Filename
	title := fileName[:len(fileName)-3]
	content, err := readFromFileHeader(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	defaultImg := dao.GetConfig(db, g.CONFIG_ARTICLE_COVER)
	err = dao.ImportArticle(db, auth.ID, title, content, defaultImg)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func readFromFileHeader(file *multipart.FileHeader) (string, error) {
	open, err := file.Open()
	if err != nil {
		slog.Error("文件读取, 目标地址错误: ", err)
		return "", err
	}
	defer open.Close()
	all, err := io.ReadAll(open)
	if err != nil {
		slog.Error("文件读取失败: ", err)
		return "", err
	}
	return string(all), nil
}
