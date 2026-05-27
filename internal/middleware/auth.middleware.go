package middleware

import (
	"strings"

	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/internal/response"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/gin-gonic/gin"
)

func RequireAuth(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.JSONUnauthorized(ctx, "Akses ditolak", "Silahkan login terlebih dahulu")
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := pkg.VerifyToken(tokenString)
		if err != nil {
			response.JSONUnauthorized(ctx, "Sesi tidak valid", "Token tidak valid atau sudah kadaluarsa")
			ctx.Abort()
			return
		}

		if authRepo.IsTokenBlacklisted(ctx.Request.Context(), userID, tokenString) {
			response.JSONUnauthorized(ctx, "Sesi telah berakhir", "Silakan login kembali")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Set("token", tokenString)

		ctx.Next()
	}
}
