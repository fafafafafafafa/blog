package ginblog

import (
	"gin-blog/internal/handle"
	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/swag/example/basic/docs"
)

var (
	// 后台管理系统接口
	categoryAPI handle.Category // 分类
	tagAPI      handle.Tag      // 标签
	userAuthAPI handle.UserAuth // 用户账号
	blogInfoAPI handle.BlogInfo // 博客设置
	uploadAPI   handle.Upload   // 文件上传
	userAPI     handle.User     // 用户
	menuAPI     handle.Menu     // 菜单
)

// TODO: 前端修改 PUT 和 PATCH 请求
func RegisterHandlers(r *gin.Engine) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerBaseHandler(r)
	registerAdminHandler(r)
	// registerBlogHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	// TODO: 登录, 注册 记录日志
	base.POST("/login", userAuthAPI.Login)          // 登录
	base.POST("/register", userAuthAPI.Register)    // 注册
	base.GET("/logout", userAuthAPI.Logout)         // 退出登录
	base.POST("/report", blogInfoAPI.Report)        // 上报信息
	base.GET("/config", blogInfoAPI.GetConfigMap)   // 获取配置
	base.PATCH("/config", blogInfoAPI.UpdateConfig) // 更新配置
	base.GET("/code", userAuthAPI.SendCode)         // 验证码
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	auth := r.Group("/api")

	// !注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())
	auth.Use(middleware.PermissionCheck())
	auth.Use(middleware.OperationLog())
	auth.Use(middleware.ListenOnline())

	auth.GET("/home", blogInfoAPI.GetHomeInfo) // 后台首页信息
	auth.POST("/upload", uploadAPI.UploadFile) // 文件上传

	// 博客设置
	setting := auth.Group("/setting")
	{
		setting.GET("/about", blogInfoAPI.GetAbout)    // 获取关于我
		setting.PUT("/about", blogInfoAPI.UpdateAbout) // 编辑关于我
	}

	// 用户模块
	user := auth.Group("/user")
	{
		user.GET("/list", userAPI.GetList)          // 用户列表
		user.PUT("", userAPI.Update)                // 更新用户信息
		user.PUT("/disable", userAPI.UpdateDisable) // 修改用户禁用状态
		// user.PUT("/password", userAPI.UpdatePassword)                // 修改普通用户密码
		user.PUT("/current/password", userAPI.UpdateCurrentPassword) // 修改管理员密码
		user.GET("/info", userAPI.GetInfo)                           // 获取当前用户信息
		user.PUT("/current", userAPI.UpdateCurrent)                  // 修改当前用户信息
		user.GET("/online", userAPI.GetOnlineList)                   // 获取在线用户
		user.POST("/offline/:id", userAPI.ForceOffline)              // 强制用户下线
	}
	// 分类模块
	category := auth.Group("/category")
	{
		category.GET("/list", categoryAPI.GetList)     // 分类列表
		category.POST("", categoryAPI.AddOrUpdate)     // 新增/编辑分类
		category.DELETE("", categoryAPI.Delete)        // 删除分类
		category.GET("/option", categoryAPI.GetOption) // 分类选项列表
	}

	//  标签模块
	tag := auth.Group("/tag")
	{
		tag.GET("/list", tagAPI.GetList)     // 标签列表
		tag.POST("", tagAPI.AddOrUpdate)     // 新增/编辑标签
		tag.DELETE("", tagAPI.Delete)        // 删除标签
		tag.GET("/option", tagAPI.GetOption) // 标签选项列表
	}
	// 文章模块
	// articles := auth.Group("/article")
	// {
	// 	articles.GET("/list", articleAPI.GetList)                 // 文章列表
	// 	articles.POST("", articleAPI.SaveOrUpdate)                // 新增/编辑文章
	// 	articles.PUT("/top", articleAPI.UpdateTop)                // 更新文章置顶
	// 	articles.GET("/:id", articleAPI.GetDetail)                // 文章详情
	// 	articles.PUT("/soft-delete", articleAPI.UpdateSoftDelete) // 软删除文章
	// 	articles.DELETE("", articleAPI.Delete)                    // 物理删除文章
	// 	articles.POST("/export", articleAPI.Export)               // 导出文章
	// 	articles.POST("/import", articleAPI.Import)               // 导入文章
	// }
	// 菜单模块
	menu := auth.Group("/menu")
	{
		// menu.GET("/list", menuAPI.GetTreeList)      // 菜单列表
		// menu.POST("", menuAPI.SaveOrUpdate)         // 新增/编辑菜单
		// menu.DELETE("/:id", menuAPI.Delete)         // 删除菜单
		menu.GET("/user/list", menuAPI.GetUserMenu) // 获取当前用户的菜单
		// menu.GET("/option", menuAPI.GetOption)      // 菜单选项列表(树形)
	}
}
