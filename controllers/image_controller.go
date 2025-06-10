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

type ImageController struct {
	config config.StaticConfig
	repo   repositories.FileRepository
}

func NewImageController(config config.StaticConfig, fileRepository repositories.FileRepository) ImageController {
	return ImageController{
		config: config,
		repo:   fileRepository,
	}
}

func (controller ImageController) GetImage(ctx *gin.Context) {
	imageUUID := strings.TrimSpace(ctx.Param("image"))
	if len(imageUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	file := controller.repo.GetFile(imageUUID, models.Image)
	if file == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.File(controller.config.ImagesDir + "/" + file)
}

func (controller ImageController) UploadImage(ctx *gin.Context) {
	file, err := controller.checkImageFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	uploadedFilename, err := controller.repo.UploadFile(ctx, file, controller.config.ImagesDir, models.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Server error, try again later"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(uploadedFilename))
}

func (controller ImageController) DeleteImage(ctx *gin.Context) {
	imageUUID := strings.TrimSpace(ctx.Param("image"))
	if len(imageUUID) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := controller.repo.DeleteFile(controller.config.ImagesDir, imageUUID, models.Image); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Image was successfully deleted"))
}

func (controller ImageController) checkImageFile(ctx *gin.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	validExtensions := []string{".jpg", ".jpeg"}
	if !slices.Contains(validExtensions, strings.ToLower(filepath.Ext(file.Filename))) {
		return nil, errors.New("valid file extensions are only .jpg or .jpeg")
	}

	if file.Size > 10_000_000 {
		return nil, errors.New("valid file size less than 10MB")
	}

	return file, nil
}
