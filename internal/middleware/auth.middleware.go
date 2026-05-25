package middleware

import (
	"strings"

	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/internal/response"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/gin-gonic/gin"
)

// RequireAuth adalah middleware untuk mengecek JWT Token
func RequireAuth(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.JSONUnauthorized(c, "Akses ditolak", "Silakan login terlebih dahulu")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.JSONUnauthorized(c, "Format token tidak valid", "Gunakan format: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		if authRepo.IsTokenBlacklisted(c.Request.Context(), tokenString) {
			response.JSONUnauthorized(c, "Sesi telah berakhir", "Silakan login kembali")
			c.Abort()
			return
		}

		userID, err := pkg.VerifyToken(tokenString)
		if err != nil {
			response.JSONUnauthorized(c, "Sesi tidak valid", "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
