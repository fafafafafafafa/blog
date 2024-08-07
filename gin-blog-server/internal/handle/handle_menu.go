package handle

import (
	"errors"
	"gin-blog/internal/dao"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct{}

type MenuTreeVO struct {
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

// 获取当前用户菜单: 生成后台管理界面的菜单
func (*Menu) GetUserMenu(c *gin.Context) {
	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	var menus []model.Menu
	var err error

	if auth.IsSuper {
		menus, err = dao.GetAllMenuList(db)
	} else {
		menus, err = dao.GetMenuListByUserId(db, auth.ID)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menus))
}

// 构建菜单列表的树形结构数据, []Menu => []MenuVo
func menus2MenuVos(menus []model.Menu) []MenuTreeVO {
	result := make([]MenuTreeVO, 0)

	firstLevelMenus := getFirstLevelMenus(menus)
	childrenMap := getMenuChildrenMap(menus)

	for _, first := range firstLevelMenus {
		menu := MenuTreeVO{Menu: first}
		for _, childMenu := range childrenMap[first.ID] {
			menu.Children = append(menu.Children, MenuTreeVO{Menu: childMenu})
		}
		delete(childrenMap, first.ID)
		result = append(result, menu)
	}

	sortMenu(result)
	return result
}

// 筛选出一级菜单 (parentId == 0 的菜单)
func getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu) // 只添加两级的菜单吗

		}
	}
	return firstLevelMenus
}

// key 是菜单 ID, value 是该菜单对应的子菜单列表
func getMenuChildrenMap(menus []model.Menu) map[int][]model.Menu {
	childrenMap := make(map[int][]model.Menu)
	for _, menu := range menus {
		if menu.ParentId != 0 {
			childrenMap[menu.ParentId] = append(childrenMap[menu.ParentId], menu)
		}
	}
	return childrenMap
}

// 以 orderNum 降序排序，包括子菜单
func sortMenu(menus []MenuTreeVO) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})
	for i := range menus {
		sort.Slice(menus[i].Children, func(j, k int) bool {
			return menus[i].Children[j].OrderNum < menus[i].Children[k].OrderNum
		})
	}
}

func (*Menu) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	menuList, _, err := dao.GetMenuList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menuList))
}

func (*Menu) AddOrUpdate(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := dao.AddOrUpdateMenu(GetDB(c), &req); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func (*Menu) Delete(c *gin.Context) {
	menuId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查要删除的菜单是否被角色使用
	use, _ := dao.CheckMenuInUse(db, menuId)
	if use {
		ReturnError(c, g.ErrMenuUsedByRole, nil)
		return
	}

	// 如果是一级菜单, 检查其是否有子菜单
	menu, err := dao.GetMenuById(db, menuId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrMenuNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 一级菜单下有子菜单, 不允许删除
	if menu.ParentId == 0 {
		has, _ := dao.CheckMenuHasChild(db, menuId)
		if has {
			ReturnError(c, g.ErrMenuHasChildren, nil)
			return
		}
	}

	if err = dao.DeleteMenu(db, menuId); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func (*Menu) GetOption(c *gin.Context) {
	menus, _, err := dao.GetMenuList(GetDB(c), "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	result := make([]model.TreeOptionVO, 0)
	for _, menu := range menus2MenuVos(menus) {
		option := model.TreeOptionVO{ID: menu.ID, Label: menu.Name}
		for _, child := range menu.Children {
			option.Children = append(option.Children, model.TreeOptionVO{ID: child.ID, Label: child.Name})
		}
		result = append(result, option)
	}

	ReturnSuccess(c, result)
}
