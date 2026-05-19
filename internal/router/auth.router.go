package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes
func RegisterAuthRoutes(rg *gin.RouterGroup, authController *controller.AuthController) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/", authController.Login)
}
