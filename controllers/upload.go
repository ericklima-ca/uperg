package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type uploadController struct{}

func NewUploadController() *uploadController {
	return &uploadController{}
}

func (uc *uploadController) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"body": gin.H{
				"error": "files not uploaded",
			},
		})
		return
	}
	files := form.File["upload[]"]
	for _, file := range files {
		c.SaveUploadedFile(file, "../tmp/uploads/"+file.Filename)
	}
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"body": gin.H{
			"msg": fmt.Sprintf("%v files uploaded", len(files)),
		},
	})
}
