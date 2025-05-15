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

type ChapterController struct {
	config config.StaticConfig
	repo   repositories.FileRepository
}

func NewChapterController(
	config config.StaticConfig,
	fileRepository repositories.FileRepository,
) ChapterController {
	return ChapterController{
		config: config,
		repo:   fileRepository,
	}
}

func (controller ChapterController) GetChapter(ctx *gin.Context) {
	chapterUUID := strings.TrimSpace(ctx.Param("chapter"))
	if len(chapterUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	file := controller.repo.GetFile(chapterUUID, models.Chapter)
	if file == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.File(controller.config.FilesDir + "/" + file)
}

func (controller ChapterController) UploadChapter(ctx *gin.Context) {
	file, err := controller.checkChapterFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	uploadedFilename, err := controller.repo.UploadFile(ctx, file, controller.config.FilesDir, models.Chapter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Server error, try again later"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(uploadedFilename))
}

func (controller ChapterController) DeleteChapter(ctx *gin.Context) {
	chapterUUID := strings.TrimSpace(ctx.Param("chapter"))
	if len(chapterUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := controller.repo.DeleteFile(controller.config.FilesDir, chapterUUID, models.Chapter); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Chapter was successfully deleted"))
}

func (controller ChapterController) checkChapterFile(ctx *gin.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	validExtensions := []string{".txt"}
	if !slices.Contains(validExtensions, strings.ToLower(filepath.Ext(file.Filename))) {
		return nil, errors.New("valid file extension is .txt")
	}

	if file.Size > 1_000_000 {
		return nil, errors.New("valid file size less than 1MB")
	}

	return file, nil
}
