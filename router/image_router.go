package router

import (
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func imageRouter(globalRoute *gin.Engine, controller controllers.ImageController) {
	api := globalRoute.Group("/api")

	imageApi := api.Group("/images")

	globalRoute.GET("/images/:image", controller.GetImage)
	imageApi.POST("/", controller.UploadImage)
	imageApi.DELETE("/:image", controller.DeleteImage)
}
