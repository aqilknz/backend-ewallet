package response

import (
	"net/http"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/gin-gonic/gin"
)

// Status 500 - Internal Server Error
func JSONInternalServerError(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusInternalServerError, dto.Response{
		Success: false,
		Message: "Error",
		Error:   err,
	})
}

// Status 400 - Bad Request
func JSONBadRequest(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusBadRequest, dto.Response{
		Success: false,
		Message: "Invalid Request Payload",
		Error:   err,
	})
}

// Status 401 - Unauthorized
func JSONUnauthorized(ctx *gin.Context, message string, err string) {
	ctx.JSON(http.StatusUnauthorized, dto.Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Status 200 - OK
func JSONSuccess(ctx *gin.Context, data any, message string) {
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Status 201 - Created
func JSONCreated(ctx *gin.Context, data any, message string) {
	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
