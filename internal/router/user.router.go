package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes mendaftarkan rute yang butuh proteksi token
func RegisterUserRoutes(rg *gin.RouterGroup, userController *controller.UserController) {
	userGroup := rg.Group("/users")

	//Semua rute di bawah userGroup otomatis harus bawa Token JWT
	userGroup.Use(middleware.RequireAuth())
	// Ambil Data
	userGroup.GET("/profile", userController.GetProfile)
	userGroup.GET("/dashboard", userController.GetDashboard)

	// Update Data
	userGroup.PUT("/profile", userController.EditProfile)
	userGroup.PUT("/password", userController.EditPassword)
	userGroup.PATCH("/pin", userController.EditPin)
}
