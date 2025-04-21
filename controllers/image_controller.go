package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
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
	err := ctx.ShouldBindJSON(&imageRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

}

func (controller ImageController) UploadBookCover(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	bookId, err := strconv.ParseUint(ctx.Param("book_id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: "\"book_id\" param was not provided",
		})
		ctx.Abort()
		return
	}

	fileExtnesion := file.Filename[strings.LastIndex(file.Filename, "."):]
	savePath := controller.Config.BookPath + "/book_" + fmt.Sprintf("%d", bookId) + fileExtnesion

	result, err := controller.uploadImage(savePath, file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	imageURL, err := controller.getImageURL(result)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Data:  imageURL,
		Error: "",
	})
	ctx.Abort()
}

func (controller ImageController) DeleteBookCover(ctx *gin.Context) {
	var imageRequest models.ImageRequest
	err := ctx.ShouldBindJSON(&imageRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	if err := controller.deleteImage(imageRequest.Path); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Data:  "Image was successfully deleted",
		Error: "",
	})
}

func (controller ImageController) UploadUserImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: "\"user_id\" param was not provided",
		})
		ctx.Abort()
		return
	}

	fileExtnesion := file.Filename[strings.LastIndex(file.Filename, "."):]
	savePath := controller.Config.UsersPath + "/user_" + fmt.Sprintf("%d", userId) + fileExtnesion

	result, err := controller.uploadImage(savePath, file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	imageURL, err := controller.getImageURL(result)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Data:  imageURL,
		Error: "",
	})
	ctx.Abort()
}

func (controller ImageController) DeleteUserImage(ctx *gin.Context) {
	var imageRequest models.ImageRequest
	err := ctx.ShouldBindJSON(&imageRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	if err := controller.deleteImage(imageRequest.Path); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Data:  nil,
			Error: err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Data:  "Image was successfully deleted",
		Error: "",
	})
}

func (controller ImageController) uploadImage(savePath string, file *multipart.FileHeader) (string, error) {
	fileData, err := file.Open()

	if err != nil {
		return "", err
	}

	defer fileData.Close()

	result, err := controller.FileClient.Upload(
		files.NewUploadArg(savePath),
		fileData,
	)

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
