package router

import (
	_ "github.com/aqilknz/backend-ewallet/docs"
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter adalah fungsi utama yang dipanggil oleh main.go
func InitRouter(app *gin.Engine, db *pgxpool.Pool) {
	// 1. Pasang Middleware CORS Global
	app.Use(middleware.CORSMiddleware)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// untuk auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(db, authRepo)
	authController := controller.NewAuthController(authService)

	// untuk user dashboard
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// untuk user transactions
	txRepo := repository.NewTransactionRepository(db)
	txService := service.NewTransactionService(txRepo)
	txController := controller.NewTransactionController(txService)

	api := app.Group("/ewallet")

	RegisterAuthRoutes(api, authController, authRepo)

	RegisterUserRoutes(api, userController, authRepo)

	RegisterTransactionRoutes(api, txController, authRepo)
}
