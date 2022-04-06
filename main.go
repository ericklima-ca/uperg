package main

import (
	"github.com/ericklima-ca/uperg/controllers"
	"github.com/ericklima-ca/uperg/router"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}
func main() {
	uploadControler := controllers.NewUploadController()
	newRouter := router.NewRouter(uploadControler)

	s := newRouter.GetServer()
	s.Run(":" + newRouter.Options.Port)
}
