package router

import (
	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/0ero-1ne/martha-storage/middlewares"
	"github.com/0ero-1ne/martha-storage/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(config config.Config, database *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	repo := repositories.NewFileRepository(database)

	imageController := controllers.NewImageController(config.StaticConfig, repo)
	chapterController := controllers.NewChapterController(config.StaticConfig, repo)
	audioController := controllers.NewAudioController(config.StaticConfig, repo)

	imageRouter(router, imageController)
	chapterRouter(router, chapterController)
	audioRouter(router, audioController)

	return router
}
