package router

import (
	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/controllers"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/gin-gonic/gin"
)

func NewRouter(config config.DropboxConfig) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	fileClient := files.New(dropbox.Config{
		Token:    config.Token,
		LogLevel: dropbox.LogDebug,
	})

	sharingClient := sharing.New(dropbox.Config{
		Token:    config.Token,
		LogLevel: dropbox.LogDebug,
	})

	imageController := controllers.ImageController{
		Config:        config,
		FileClient:    fileClient,
		SharingClient: sharingClient,
	}

	imageRouter(api.Group("/images"), imageController)

	return router
}
