package handle

import (
	"errors"
	"gin-blog/internal/dao"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Resource struct{}

// TODO: 考虑删除模块后, 其子资源怎么办? 目前做法是有子资源无法删除
// TODO: 强制删除?
func (*Resource) Delete(c *gin.Context) {
	resourceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查该资源是否被角色使用
	use, _ := dao.CheckResourceInUse(db, resourceId)
	if use {
		ReturnError(c, g.ErrResourceUsedByRole, nil)
		return
	}

	// 获取该资源
	resource, err := dao.GetResourceById(db, resourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrResourceNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 如果作为模块, 检查模块下是否有子资源
	if resource.ParentId == 0 {
		hasChild, _ := dao.CheckResourceHasChild(db, resourceId)
		if hasChild {
			ReturnError(c, g.ErrResourceHasChildren, nil)
			return
		}
	}

	rows, err := dao.DeleteResource(db, resourceId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// 新增或编辑资源, 关联更新 casbin_rule 中数据
func (*Resource) AddOrUpdate(c *gin.Context) {
	var req model.AddOrEditResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	err := dao.AddOrUpdateResource(db, req.ID, req.ParentId, req.Name, req.Url, req.Method)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 编辑资源的匿名访问, 关联更新 casbin_rule 中数据
func (*Resource) UpdateAnonymous(c *gin.Context) {
	var req model.EditAnonymousReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := dao.UpdateResourceAnonymous(GetDB(c), req.ID, req.Anonymous)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 获取资源列表(树形)
func (*Resource) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	resourceList, err := dao.GetResourceList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, resources2ResourceVos(resourceList))
}

// []Resource => []ResourceVO
func resources2ResourceVos(resources []model.Resource) []model.ResourceTreeVO {
	list := make([]model.ResourceTreeVO, 0)
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	for _, item := range parentList {
		resourceVO := resource2ResourceVo(item)
		resourceVO.Children = make([]model.ResourceTreeVO, 0)
		for _, child := range childrenMap[item.ID] {
			resourceVO.Children = append(resourceVO.Children, resource2ResourceVo(child))
		}
		list = append(list, resourceVO)
	}
	return list
}

// Resource => ResourceVO
func resource2ResourceVo(r model.Resource) model.ResourceTreeVO {
	return model.ResourceTreeVO{
		ID:        r.ID,
		Name:      r.Name,
		Url:       r.Url,
		Method:    r.Method,
		Anonymous: r.Anonymous,
		CreatedAt: r.CreatedAt,
	}
}

// 获取一级资源 (parent_id == 0)
func getModuleList(resources []model.Resource) []model.Resource {
	list := make([]model.Resource, 0)
	for _, r := range resources {
		if r.ParentId == 0 {
			list = append(list, r)
		}
	}
	return list
}

// 存储每个节点对应 [子资源列表] 的 map
// key: resourceId
// value: childrenList
func getChildrenMap(resources []model.Resource) map[int][]model.Resource {
	m := make(map[int][]model.Resource)
	for _, r := range resources {
		if r.ParentId != 0 {
			m[r.ParentId] = append(m[r.ParentId], r)
		}
	}
	return m
}

// 获取数据选项(树形)
func (*Resource) GetOption(c *gin.Context) {
	result := make([]model.TreeOptionVO, 0)

	db := GetDB(c)
	resources, err := dao.GetResourceList(db, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)

	for _, item := range parentList {
		var children []model.TreeOptionVO
		for _, re := range childrenMap[item.ID] {
			children = append(children, model.TreeOptionVO{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		result = append(result, model.TreeOptionVO{
			ID:       item.ID,
			Label:    item.Name,
			Children: children,
		})
	}
	ReturnSuccess(c, result)
}
