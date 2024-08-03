package routers

import (
	"go-gin-example/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // 其实就是gin.Default()

	gin.SetMode(setting.ServerSetting.RunMode)

	// 无需鉴权的接口
	base := r.Group("/api")
	{
		// TODO: 用户注册 和 后台登录 应该记录到 日志
		// base.POST("/login", userAuthAPI.Login) // 后台登录]
		base.POST("/login", GetAuth) // 后台登录]
		// base.POST("/report", blogInfoAPI.Report) // 上报信息
	}

	// r.GET("/auth", GetAuth)
	// r.POST("/upload", api.UploadImage)
	// r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	// apiv1 := r.Group("/api/v1")
	// apiv1.Use(jwt.JWT())
	// {
	// 	apiv1.GET("/tags", v1.GetTags)
	// 	apiv1.POST("/tags", v1.AddTag)
	// 	apiv1.PUT("/tags/:id", v1.EditTag)
	// 	apiv1.DELETE("/tags/:id", v1.DeleteTag)

	// 	apiv1.GET("/articles", v1.GetArticles)
	// 	apiv1.GET("/articles/:id", v1.GetArticle)
	// 	apiv1.POST("/articles", v1.AddArticle)
	// 	apiv1.PUT("/articles/:id", v1.EditArticle)
	// 	apiv1.DELETE("/articles/:id", v1.DeleteArticle)

	// }

	return r
}
