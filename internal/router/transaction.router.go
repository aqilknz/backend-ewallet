package router

import (
	"github.com/aqilknz/backend-ewallet/internal/controller"
	"github.com/aqilknz/backend-ewallet/internal/middleware"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup, txController *controller.TransactionController, authRepo repository.AuthRepository) {
	txGroup := rg.Group("/users/transaction")
	txGroup.Use(middleware.RequireAuth(authRepo))

	txGroup.POST("/topup", txController.TopUp)
	txGroup.POST("/transfer", txController.Transfer)
	txGroup.GET("/history", txController.GetHistory)
	txGroup.GET("/report", txController.GetReport)
}
