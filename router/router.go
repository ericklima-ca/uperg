package router

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	UploadHandler Uploader
	Options       ServerOptions
}

type ServerOptions struct {
	Port string
}

type Uploader interface {
	UploadFiles(*gin.Context)
}

func NewRouter(uploader Uploader) *Router {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	return &Router{
		UploadHandler: uploader,
		Options: ServerOptions{
			Port: port,
		},
	}
}

func (r *Router) GetServer() *gin.Engine {
	engine := gin.Default()

	uploadAPI := engine.Group("api/")
	uploadAPI.POST("/api/uploads", r.UploadHandler.UploadFiles)

	return engine
}
