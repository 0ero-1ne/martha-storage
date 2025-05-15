package router

import (
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func audioRouter(globalRoute *gin.Engine, controller controllers.AudioController) {
	api := globalRoute.Group("/api")

	audioApi := api.Group("/audios")

	globalRoute.GET("/audios/:audio", controller.GetAudio)
	audioApi.POST("", controller.UploadAudio)
	audioApi.DELETE("/:audio", controller.DeleteAudio)
}
