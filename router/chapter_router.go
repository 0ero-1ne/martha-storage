package router

import (
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func chapterRouter(globalRoute *gin.Engine, controller controllers.ChapterController) {
	api := globalRoute.Group("/api")

	chapterApi := api.Group("/chapters")

	globalRoute.GET("/chapters/:chapter", controller.GetChapter)
	chapterApi.POST("", controller.UploadChapter)
	chapterApi.DELETE("/:chapter", controller.DeleteChapter)
}
