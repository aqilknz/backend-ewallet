package middleware

import (
	"net/http"
	"strings"

	"github.com/aqilknz/backend-ewallet/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RequireAuth adalah middleware untuk mengecek JWT Token
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Akses ditolak, silakan login terlebih dahulu",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Format token tidak valid (Gunakan: Bearer <token>)",
			})
			return
		}

		tokenString := parts[1]

		userID, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau sudah kadaluarsa",
			})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
