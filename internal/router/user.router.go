package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes mendaftarkan rute yang butuh proteksi token
func RegisterUserRoutes(rg *gin.RouterGroup, userController *controller.UserController, authRepo repository.AuthRepository) {
	userGroup := rg.Group("/users")

	//Semua rute di bawah userGroup otomatis harus bawa Token JWT
	userGroup.Use(middleware.RequireAuth(authRepo))
	// Ambil Data
	userGroup.GET("/profile", userController.GetProfile)
	userGroup.GET("/dashboard", userController.GetDashboard)

	// Update Data
	userGroup.PATCH("/profile", userController.EditProfile)
	userGroup.PATCH("/profile/password", userController.EditPassword)
	userGroup.PATCH("/profile/pin", userController.EditPin)
	userGroup.POST("/transaction/checkpin", userController.CheckPin)
	userGroup.GET("/receivers", userController.FindReceivers)
}
