package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/models"
	"github.com/0ero-1ne/martha-storage/utils"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	Config config.StaticConfig
}

func NewImageController(config config.StaticConfig) ImageController {
	return ImageController{
		Config: config,
	}
}

func (controller ImageController) UploadImage(ctx *gin.Context) {
	file, err := controller.loadFormFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	fileExt := filepath.Ext(file.Filename)

	randomizer := utils.NewRandomizer()
	randomFileName, _ := randomizer.GenerateString(16)

	savePath := fmt.Sprintf("%s/%s%s", controller.Config.ImagesDir, randomFileName, fileExt)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Server error, try again later"))
		return
	}

	returnValue := fmt.Sprintf("%s/images/%s%s", ctx.Request.Host, randomFileName, fileExt)
	ctx.JSON(http.StatusOK, models.NewSuccessResponse(returnValue))
}

func (controller ImageController) DeleteImage(ctx *gin.Context) {
	imageName := ctx.Param("image")

	if len(imageName) == 0 {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid image param value"))
		return
	}

	if err := os.Remove(controller.Config.ImagesDir + "/" + imageName); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid image param value"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Image was successfully deleted"))
}

func (controller ImageController) loadFormFile(ctx *gin.Context) (*multipart.FileHeader, error) {
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
