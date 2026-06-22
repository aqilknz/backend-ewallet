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
	authGroup.POST("/create-pin", middleware.RequireAuth(authRepo), authController.CreatePin)
	authGroup.POST("", authController.Login)
	authGroup.POST("/forgot-password", authController.ForgotPassword)
	authGroup.POST("/verify-otp", authController.VerifyOTP)
	authGroup.POST("/reset-password", authController.ResetPassword)
	authGroup.DELETE("/logout", middleware.RequireAuth(authRepo), authController.Logout)
	authGroup.POST("/check-email", authController.CheckEmail)
	authGroup.POST("/update-password", authController.UpdatePassword)
}
