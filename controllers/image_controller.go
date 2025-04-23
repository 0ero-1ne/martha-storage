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
	image := strings.TrimSpace(ctx.Param("image"))
	if len(image) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	file := controller.repo.GetFile(image)
	ctx.File(controller.config.ImagesDir + "/" + file)
}

func (controller ImageController) UploadImage(ctx *gin.Context) {
	file, err := controller.checkImage(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	uploadedFilename, err := controller.repo.UploadFile(ctx, file, controller.config.ImagesDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Server error, try again later"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(uploadedFilename))
}

func (controller ImageController) DeleteImage(ctx *gin.Context) {
	image := strings.TrimSpace(ctx.Param("image"))
	if len(image) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := controller.repo.DeleteFile(controller.config.ImagesDir, image); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Image was successfully deleted"))
}

func (controller ImageController) checkImage(ctx *gin.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	validExtensions := []string{".jpg", ".jpeg"}
	if !slices.Contains(validExtensions, strings.ToLower(filepath.Ext(file.Filename))) {
		return nil, errors.New("Valid file extensions are only .jpg or .jpeg")
	}

	if file.Size > 500_000 {
		return nil, errors.New("Valid file size less than 500Kb")
	}

	return file, nil
}
