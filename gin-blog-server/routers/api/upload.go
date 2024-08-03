package api

import (
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/logging"
	"log"
	"net/http"

	"go-gin-example/pkg/upload"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	var code int
	data := make(map[string]interface{})
	code = e.SUCCESS

	file, image, err := c.Request.FormFile("image")

	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()
		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				log.Println(err)
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL

			} else if err := c.SaveUploadedFile(image, src); err != nil {
				log.Println(err)
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL

			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
