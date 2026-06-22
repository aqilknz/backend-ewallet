package router

import (
	_ "github.com/aqilknz/backend-ewallet/docs"
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter adalah fungsi utama yang dipanggil oleh main.go
func InitRouter(app *gin.Engine, db *pgxpool.Pool, redis *redis.Client) {
	// Pasang Middleware CORS Global
	app.Use(middleware.CORSMiddleware)
	app.Static("/ewallet/img/profiles", "./public/img/profiles")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// untuk auth
	authRepo := repository.NewAuthRepository(db, redis)
	authService := service.NewAuthService(db, authRepo)
	authController := controller.NewAuthController(authService)

	// untuk user dashboard
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, redis)
	userController := controller.NewUserController(userService)

	// untuk user transactions
	txRepo := repository.NewTransactionRepository(db)
	txService := service.NewTransactionService(txRepo, redis)
	txController := controller.NewTransactionController(txService)

	api := app.Group("/ewallet")

	RegisterAuthRoutes(api, authController, authRepo)

	RegisterUserRoutes(api, userController, authRepo)

	RegisterTransactionRoutes(api, txController, authRepo)
}
