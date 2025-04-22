package router

import (
	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func NewRouter(config config.Config) *gin.Engine {
	router := gin.Default()

	imageController := controllers.NewImageController(config.StaticConfig)

	imageRouter(router, imageController)

	return router
}
