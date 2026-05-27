package response

import (
	"net/http"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/gin-gonic/gin"
)

// Status 500 - Internal Server Error
func JSONInternalServerError(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
		Success: false,
		Message: "Error",
		Error:   err,
	})
}

// Status 400 - Bad Request
func JSONBadRequest(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusBadRequest, dto.Response[any]{
		Success: false,
		Message: "Invalid Request Payload",
		Error:   err,
	})
}

// Status 401 - Unauthorized
func JSONUnauthorized(ctx *gin.Context, message string, err string) {
	ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Status 200 - OK
func JSONSuccess[T any](ctx *gin.Context, data T, message string) {
	ctx.JSON(http.StatusOK, dto.Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Status 201 - Created
func JSONCreated[T any](ctx *gin.Context, data T, message string) {
	ctx.JSON(http.StatusCreated, dto.Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Status 404 - Not Found
// Digunakan ketika data yang di-request tidak ditemukan di database
func JSONNotFound(ctx *gin.Context, message string, err string) {
	ctx.JSON(http.StatusNotFound, dto.Response[any]{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Status 409 - Conflict
func JSONConflict(ctx *gin.Context, message string, err string) {
	ctx.JSON(http.StatusConflict, dto.Response[any]{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// status 422 - Unprocess
func JSONUnprocessableEntity(ctx *gin.Context, message string, err string) {
	ctx.JSON(http.StatusUnprocessableEntity, dto.Response[any]{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// status 204 - No Content
func JSONNoContent(ctx *gin.Context, message string) {
	ctx.Status(http.StatusNoContent)
}
