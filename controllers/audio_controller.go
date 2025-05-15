package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/models"
	"github.com/0ero-1ne/martha-storage/repositories"
	"github.com/gin-gonic/gin"
)

type AudioController struct {
	config config.StaticConfig
	repo   repositories.FileRepository
}

func NewAudioController(
	config config.StaticConfig,
	fileRepository repositories.FileRepository,
) AudioController {
	return AudioController{
		config: config,
		repo:   fileRepository,
	}
}

func (controller AudioController) GetAudio(ctx *gin.Context) {
	audioUUID := strings.TrimSpace(ctx.Param("audio"))
	if len(audioUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	file := controller.repo.GetFile(audioUUID, models.Audio)
	if file == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.File(controller.config.AudiosDir + "/" + file)
}

func (controller AudioController) UploadAudio(ctx *gin.Context) {
	file, err := controller.checkAudioFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	uploadedFilename, err := controller.repo.UploadFile(ctx, file, controller.config.AudiosDir, models.Audio)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Server error, try again later"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(uploadedFilename))
}

func (controller AudioController) DeleteAudio(ctx *gin.Context) {
	audioUUID := strings.TrimSpace(ctx.Param("audio"))
	if len(audioUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := controller.repo.DeleteFile(controller.config.AudiosDir, audioUUID, models.Audio); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Audio was successfully deleted"))
}

func (controller AudioController) checkAudioFile(ctx *gin.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	validExtensions := []string{".mp3", ".flac"}
	if !slices.Contains(validExtensions, strings.ToLower(filepath.Ext(file.Filename))) {
		return nil, errors.New("valid file extension is .mp3 and .flac")
	}

	if file.Size > 30_000_000 {
		return nil, errors.New("valid file size less than 30MB")
	}

	return file, nil
}
