package ginblog

import (
	"gin-blog/internal/handle"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/swag/example/basic/docs"
)

var (
	// 后台管理系统接口
	userAuthAPI handle.UserAuth // 用户账号
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

}
