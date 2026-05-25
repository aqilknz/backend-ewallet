package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware adalah fungsi penjaga untuk mengatur siapa saja yang boleh mengakses API
// "http://127.0.0.1:5500"
func CORSMiddleware(ctx *gin.Context) {
	allowedOrigin := []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:9000"}

	currentOrigin := ctx.GetHeader("Origin")
	if slices.Contains(allowedOrigin, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
	}

	allowedHeaders := []string{"Content-Type", "Authorization"}
	ctx.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}
	ctx.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
