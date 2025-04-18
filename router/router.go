package router

import (
	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func NewRouter(config config.DropboxConfig) *gin.Engine {
	router := gin.New()
	api := router.Group("/api")

	bookRoutes(api)
	userRoutes(api)

	return router
}

func bookRoutes(globalRoute *gin.RouterGroup) {
	bookRoutes := globalRoute.Group("/books")
	controller := controllers.BookController{}

	bookRoutes.POST("/upload", controller.UploadImage)
	bookRoutes.DELETE("/delete", controller.DeleteImage)
}

func userRoutes(globalRoute *gin.RouterGroup) {
	userRoutes := globalRoute.Group("/users")
	controller := controllers.UserController{}

	userRoutes.POST("/upload", controller.UploadImage)
	userRoutes.DELETE("/delete", controller.DeleteImage)
}
