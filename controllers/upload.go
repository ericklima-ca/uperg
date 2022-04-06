package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ericklima-ca/uperg/tasks"
	"github.com/gin-gonic/gin"
)

var path = "./tmp/uploads/"

type uploadController struct{}

func NewUploadController() *uploadController {
	return &uploadController{}
}

func (uc *uploadController) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"body": gin.H{
				"error": "files not uploaded",
			},
		})
		return
	}
	files := form.File["upload[]"]
	os.MkdirAll(path, 0777)
	defer os.RemoveAll(path)
	for _, file := range files {
		c.SaveUploadedFile(file, path+file.Filename)
	}
	doneProcesses := make(chan bool)

	start := time.Now()
	go func() {
		tasks.ProcessFiles(doneProcesses, path)
	}()
	
	<-doneProcesses
	log.Println("Process finished")
	
	log.Println(time.Since(start).Seconds())

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"body": gin.H{
			"msg": fmt.Sprintf("%v files uploaded", len(files)),
		},
	})
}
