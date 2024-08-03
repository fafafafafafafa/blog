package util

import (
	"go-gin-example/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {

	res := 0
	page := com.StrTo(c.Query("page")).MustInt()
	if page > 0 {
		res = (page - 1) * setting.AppSetting.PageSize
	}
	return res
}
