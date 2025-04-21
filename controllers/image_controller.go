package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/models"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	Config        config.DropboxConfig
	FileClient    files.Client
	SharingClient sharing.Client
}

func (controller ImageController) GetImageURL(ctx *gin.Context) {
	var imageRequest models.ImageRequest
	if err := ctx.ShouldBindJSON(&imageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	imageURL, err := controller.getImageURL(imageRequest.Path)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(imageURL))
}

func (controller ImageController) UploadBookCover(ctx *gin.Context) {
	file, err := controller.loadFormFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	fileExt := filepath.Ext(file.Filename)
	bookId := ctx.GetUint("book_id")
	savePath := fmt.Sprintf("%s/book_%d%s", controller.Config.BookPath, bookId, fileExt)

	result, err := controller.uploadImage(savePath, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(result))
}

func (controller ImageController) DeleteBookCover(ctx *gin.Context) {
	var imageRequest models.ImageRequest

	if err := ctx.ShouldBindJSON(&imageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	if err := controller.deleteImage(imageRequest.Path); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Image was successfully deleted"))
}

func (controller ImageController) UploadUserImage(ctx *gin.Context) {
	file, err := controller.loadFormFile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	fileExt := filepath.Ext(file.Filename)
	userId := ctx.GetUint("user_id")
	savePath := fmt.Sprintf("%s/user_%d%s", controller.Config.BookPath, userId, fileExt)

	result, err := controller.uploadImage(savePath, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(result))
}

func (controller ImageController) DeleteUserImage(ctx *gin.Context) {
	var imageRequest models.ImageRequest

	if err := ctx.ShouldBindJSON(&imageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	if err := controller.deleteImage(imageRequest.Path); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse("Image was successfully deleted"))
}

func (controller ImageController) uploadImage(savePath string, file *multipart.FileHeader) (string, error) {
	fileData, err := file.Open()

	if err != nil {
		return "", err
	}

	defer fileData.Close()

	result, err := controller.FileClient.Upload(files.NewUploadArg(savePath), fileData)

	if err != nil {
		return "", err
	}

	return result.PathDisplay, nil
}

func (controller ImageController) deleteImage(path string) error {
	_, err := controller.FileClient.DeleteV2(files.NewDeleteArg(path))

	return err
}

func (controller ImageController) getImageURL(path string) (string, error) {
	result, err := controller.SharingClient.CreateSharedLink(
		sharing.NewCreateSharedLinkArg(path),
	)

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(result.Url, "&dl=0", "&dl=1"), nil
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
