package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup, txController *controller.TransactionController, authRepo repository.AuthRepository) {
	// userGroup := rg.Group("/users")
	txGroup := rg.Group("/users")
	txGroup.Use(middleware.RequireAuth(authRepo))

	txGroup.POST("/transaction/topup", txController.TopUp)
	txGroup.POST("/transaction/transfer", txController.Transfer)
	txGroup.GET("/transaction/history", txController.GetHistory)
	txGroup.GET("/transaction/report", txController.GetReport)
}
