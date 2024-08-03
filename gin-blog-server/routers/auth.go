package routers

import (
	"go-gin-example/pkg/app"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/util"
	"go-gin-example/service/auth_services"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary      Get Auth
// @Produce      json
// @Param        username	query	string	true  "username"
// @Param        password	query	string	true  "password"
// @Success 	 200 {string} json "{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{
		Username: username,
		Password: password,
	}
	ok, _ := valid.Valid(&a)

	appG := app.Gin{C: c}
	if !ok {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	serviceAuth := auth_services.Auth{
		Username: username,
		Password: password,
	}
	exist, err := serviceAuth.CheckAuth()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	data := make(map[string]interface{})
	data["token"] = token
	appG.Response(http.StatusOK, e.SUCCESS, data)

}
