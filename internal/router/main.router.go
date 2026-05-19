package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/aqilknz/backend-ewallet/pkg/utils" // Import CORS kamu
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitRouter adalah fungsi utama yang dipanggil oleh main.go
func InitRouter(app *gin.Engine, db *pgxpool.Pool) {
	// 1. Pasang Middleware CORS Global
	app.Use(utils.CORSMiddleware)

	// untuk auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(db, authRepo)
	authController := controller.NewAuthController(authService)

	// untuk user dashboard
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	api := app.Group("/ewallet")

	RegisterAuthRoutes(api, authController)

	RegisterUserRoutes(api, userController)
}
