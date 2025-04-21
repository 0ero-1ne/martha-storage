package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/0ero-1ne/martha-storage/models"
	"github.com/gin-gonic/gin"
)

func ParseParamsId(params []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, param := range params {
			id, err := strconv.ParseUint(ctx.Param(param), 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(fmt.Sprintf("Invalid '%s' param value", param)))
				return
			}
			ctx.Set(param, uint(id))
		}
		ctx.Next()
	}
}
