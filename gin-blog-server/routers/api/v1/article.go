package v1

import (
	"go-gin-example/pkg/app"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"go-gin-example/service/article_service"
	"go-gin-example/service/tag_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取文章列表
// @Summary      Get Articles
// @Produce      json
// @Param        state	query	int	false  "State"
// @Param        tag_id	query	int	false  "TagId"
// @Success 	 200 {string} json "{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")

	}
	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()
		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")

	}
	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	serviceArticle := article_service.Article{
		TagId:    tag_id,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := serviceArticle.GetArticlesTotal()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return

	}
	list, err := serviceArticle.GetArticles()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["list"] = list
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// 获取指定文章
// @Summary      Get Single Article By Id
// @Produce      json
// @Param        id						query	int		true	"ID"
// @Success 	 200 {string} 			json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} 	[get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章Id 最小为1")

	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	serviceArticle := article_service.Article{
		ID: id,
	}
	exist, err := serviceArticle.ExistArticleByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	var data interface{}

	data, err = serviceArticle.GetArticle()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// 新建文章
// @Summary      Add Article
// @Produce      json
// @Param        title			query	string	true  	"Title"
// @Param        desc			query	string	true  	"Desc"
// @Param        content		query	string	true  	"Content"
// @Param        created_by		query	string	true  	"CreatedBy"
// @Param        tag_id			query	int		false  	"TagId"
// @Success 	 200 {string} 	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	valid := validation.Validation{}

	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	created_by := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	cover_image_url := c.Query("cover_image_url")

	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()

		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")
	}

	valid.Required(title, "title").Message("文章标题不为空")
	valid.MaxSize(title, 100, "title").Message("文章标题最大长度为100字符")
	valid.Required(desc, "desc").Message("文章简介不为空")
	valid.MaxSize(desc, 255, "desc").Message("文章简介最大长度为255字符")
	valid.Required(content, "content").Message("文章内容不为空")
	valid.MaxSize(content, 65535, "content").Message("文章内容最大长度为65535字符")
	valid.Required(created_by, "created_by").Message("文章创建人不为空")
	valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")
	valid.Required(cover_image_url, "cover_image_url").Message("文章封面url不为空")
	valid.MaxSize(cover_image_url, 255, "cover_image_url").Message("文章封面url最大长度为255字符")

	appG := app.Gin{C: c}
	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	serviceTag := tag_service.Tag{ID: tag_id}
	exist, err := serviceTag.ExistTagByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	serviceArticle := article_service.Article{
		TagId:         tag_id,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CreatedBy:     created_by,
		State:         state,
		CoverImageUrl: cover_image_url,
	}
	err = serviceArticle.AddArticle()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 更新指定文章
// @Summary      Update Article By Id
// @Produce      json
// @Param        id				query		int		true  	"ID"
// @Param        title			query		string	true 	"Title"
// @Param        desc			query		string	true  	"Desc"
// @Param        content		query		string	true  	"Content"
// @Param        modified_by	query		string	true  	"ModifiedBy"
// @Param        tag_id			query		int		false  	"TagId"
// @Param        state			query		int		false  	"State"
// @Success 	 200 			{string}	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modified_by := c.Query("modified_by")
	cover_image_url := c.Query("cover_image_url")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")
	}
	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()

		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")
	}

	valid.Min(id, 1, "id").Message("文章id 最小为1")
	valid.Required(title, "title").Message("文章标题不为空")
	valid.MaxSize(title, 100, "title").Message("文章标题最大长度为100字符")
	valid.Required(desc, "desc").Message("文章简介不为空")
	valid.MaxSize(desc, 255, "desc").Message("文章简介最大长度为255字符")
	valid.Required(content, "content").Message("文章内容不为空")
	valid.MaxSize(content, 65535, "content").Message("文章内容最大长度为65535字符")
	valid.Required(modified_by, "modified_by").Message("文章更改人不为空")
	valid.MaxSize(modified_by, 100, "modified_by").Message("文章更改人最大长度为100字符")
	valid.Required(cover_image_url, "cover_image_url").Message("文章封面url不为空")
	valid.MaxSize(cover_image_url, 100, "cover_image_url").Message("文章封面url最大长度为100字符")

	appG := app.Gin{C: c}
	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	serviceTag := tag_service.Tag{
		ID: tag_id,
	}
	exist, err := serviceTag.ExistTagByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	serviceArticle := article_service.Article{
		ID:            id,
		TagId:         tag_id,
		State:         state,
		Title:         title,
		Desc:          desc,
		Content:       content,
		ModifiedBy:    modified_by,
		CoverImageUrl: cover_image_url,
	}
	exist, err = serviceArticle.ExistArticleByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	err = serviceArticle.EditArticle()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 删除指定文章
// @Summary      Delete Article By Id
// @Produce      json
// @Param        id		query		int		true  "ID"
// @Success 	 200 	{string} 	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} 	[delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章Id 最小为1")

	appG := app.Gin{C: c}
	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	serviceArticle := article_service.Article{ID: id}
	exist, err := serviceArticle.ExistArticleByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	err = serviceArticle.DeleteArticle()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
