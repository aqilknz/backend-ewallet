package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes
func RegisterAuthRoutes(rg *gin.RouterGroup, authController *controller.AuthController, authRepo repository.AuthRepository) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/", authController.Login)
	authGroup.DELETE("/logout", middleware.RequireAuth(authRepo), authController.Logout)
}
