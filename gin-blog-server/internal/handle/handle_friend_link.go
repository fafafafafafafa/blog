package handle

import (
	"gin-blog/internal/dao"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Link struct{}

// @Summary 获取友链列表
// @Description 根据条件查询获取友链列表
// @Tags Link
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Param keyword query string false "搜索关键字"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.FriendLink]]
// @Security ApiKeyAuth
// @Router /link/list [get]
func (*Link) GetList(c *gin.Context) {
	var query model.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := dao.GetLinkList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, model.PageResult[model.FriendLink]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 添加或修改友链
// @Description 添加或修改友链
// @Tags Link
// @Param form body AddOrEditLinkReq true "添加或修改友链"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.FriendLink]
// @Security ApiKeyAuth
// @Router /link [post]
func (*Link) AddOrUpdate(c *gin.Context) {
	var req model.AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	link, err := dao.AddOrUpdateLink(GetDB(c), req.ID, req.Name, req.Avatar, req.Address, req.Intro)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, link)
}

// @Summary 删除友链（批量）
// @Description 根据 ID 数组删除友链
// @Tags Link
// @Param ids body []int true "友链ID数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /link [delete]
func (*Link) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(&model.FriendLink{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
