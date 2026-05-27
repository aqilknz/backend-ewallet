package controller

import (
	"errors"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/response"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	txService *service.TransactionService
}

func NewTransactionController(txService *service.TransactionService) *TransactionController {
	return &TransactionController{txService: txService}
}

// Top Up Saldo
//
//	@Summary		User Top Up
//	@Description	Add balance to user's wallet
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body		dto.TopUpRequest true "Top Up Payload"
//	@Success		200 {object}	dto.Response[dto.TopUpResponse]
//	@Failure		422 {object}	dto.Response[any]
//	@Failure		401 {object}	dto.Response[any]
//	@Failure		500 {object}	dto.Response[any]
//	@Router			/users/transaction/topup [post]
func (tc *TransactionController) TopUp(ctx *gin.Context) {
	var req dto.TopUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONUnprocessableEntity(ctx, "Pastikan minimal top up Rp 10.000.", "Data input tidak valid.")
		return
	}

	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		response.JSONUnauthorized(ctx, "Akses ditolak", "Token tidak valid")
		return
	}
	userID := userIDRaw.(int)

	res, err := tc.txService.TopUp(ctx.Request.Context(), userID, req)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, res, "Top Up berhasil diproses")
}

// Transfer Saldo
//
//	@Summary		User Transfer
//	@Description	transfer saldo to receiver with pin
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body		dto.TransferRequest true "Transfer Payload"
//	@Success		200 {object}	dto.Response[dto.TransferResponse]
//	@Failure		400 {object}	dto.Response[any]
//	@Failure		401 {object}	dto.Response[any]
//	@Failure		422 {object}	dto.Response[any]
//	@Failure		500 {object}	dto.Response[any]
//	@Router			/users/transaction/transfer [post]
func (tc *TransactionController) Transfer(ctx *gin.Context) {
	var req dto.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Format input tidak sesuai. Pastikan nominal minimal Rp 10.000 dan PIN terisi 6 digit angka.")
		return
	}

	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		response.JSONUnauthorized(ctx, "Akses ditolak", "Token tidak valid")
		return
	}
	senderID := userIDRaw.(int)

	res, err := tc.txService.Transfer(ctx.Request.Context(), senderID, req)
	if err != nil {
		if errors.Is(err, service.ErrPinNotSet) {
			response.JSONUnprocessableEntity(ctx, "Validasi akun gagal", err.Error())
			return
		}

		response.JSONUnprocessableEntity(ctx, "Tidak bisa memproses transaksi", err.Error())
		return
	}

	response.JSONSuccess(ctx, res, "Transfer berhasil diproses")
}

// Get Transaction History
//
//	@Summary        Get transaction history
//	@Description    Retrieve user's transaction history (topup, transfer) with search and pagination
//	@Tags           transaction
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          search  query   string  false   "Search in description"
//	@Param          page    query   int     false   "Page number"       default(1)
//	@Param          limit   query   int     false   "Items per page"    default(10)
//	@Success        200     {object}    dto.Response[dto.TransactionHistoryResponse]
//	@Failure        400     {object}    dto.Response[any]
//	@Failure        500     {object}    dto.Response[any]
//	@Router         /users/transaction/history [get]
func (tc *TransactionController) GetHistory(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var param dto.TransactionHistoryFilterParam

	if err := ctx.ShouldBindQuery(&param); err != nil {
		response.JSONBadRequest(ctx, "Parameter query tidak valid")
		return
	}

	result, err := tc.txService.GetHistory(ctx.Request.Context(), userID, param)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, result, "Berhasil mengambil riwayat transaksi")
}

// Get Transaction Report
//
//	@Summary        Get transaction report for graph
//	@Description    Retrieve aggregated transaction data grouped by date for dashboard charts
//	@Tags           transaction
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          type        query   string  false   "Filter by type (income, expense, both)" Enums(income, expense, both) default(both)
//	@Param          start_date  query   string  false   "Start date (YYYY-MM-DD) e.g., 2026-05-01"
//	@Param          end_date    query   string  false   "End date (YYYY-MM-DD) e.g., 2026-05-31"
//	@Success        200         {object}    dto.Response[[]dto.TransactionReportItem]
//	@Failure        400         {object}    dto.Response[any]
//	@Failure        500         {object}    dto.Response[any]
//	@Router         /users/transaction/report [get]
func (tc *TransactionController) GetReport(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var param dto.TransactionReportFilterParam

	if err := ctx.ShouldBindQuery(&param); err != nil {
		response.JSONBadRequest(ctx, "Parameter query tidak valid")
		return
	}

	result, err := tc.txService.GetReport(ctx.Request.Context(), userID, param)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, result, "Berhasil mengambil data laporan grafik transaksi")
}
