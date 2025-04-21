package router

import (
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/gin-gonic/gin"
)

func imageRouter(globalRoute *gin.RouterGroup, controller controllers.ImageController) {
	globalRoute.GET("/", controller.GetImageURL)

	bookRoutes := globalRoute.Group("/books")
	bookRoutes.POST("/:book_id", controller.UploadBookCover)
	bookRoutes.DELETE("/", controller.DeleteBookCover)

	userRoutes := globalRoute.Group("/users")
	userRoutes.POST("/:user_id", controller.UploadUserImage)
	userRoutes.DELETE("/", controller.DeleteUserImage)
}
