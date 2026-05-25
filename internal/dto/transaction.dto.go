package dto

import "time"

type CheckPinRequest struct {
	Pin string `json: pin binding:"required, len=6, numeric"`
}

// ==========================================
// DTO TOPUP
// ==========================================
type TopUpRequest struct {
	Amount          int `json:"amount" binding:"required,gt=0"`
	PaymentMethodID int `json:"payment_method_id" binding:"required"`
	Discount        int `json:"discount"`
	Tax             int `json:"tax"`
	SubTotal        int `json:"sub_total" binding:"required"`
}

type TopUpResponse struct {
	TransactionID   int       `json:"transaction_id"`
	Amount          int       `json:"amount"`
	PaymentMethodID int       `json:"payment_method_id"`
	Discount        int       `json:"discount"`
	Tax             int       `json:"tax"`
	SubTotal        int       `json:"sub_total"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// ==========================================
// DTO TRANSFER
// ==========================================
type TransferRequest struct {
	ReceiverEmail string `json:"receiver_email" binding:"required,email"`
	Amount        int    `json:"amount" binding:"required,gt=0"`
	Notes         string `json:"notes"`
}

type TransferResponse struct {
	TransactionID  int       `json:"transaction_id"`
	SenderID       int       `json:"sender_id"`
	CounterpartyID int       `json:"counterparty_id"`
	Amount         int       `json:"amount"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

// ==========================================
// DTO HISTORY & REPORT
// ==========================================
type TransactionHistoryFilterParam struct {
	Search string `form:"search"`
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type TransactionHistoryItem struct {
	ID              int       `json:"id"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"` // topup / transfer
	FlowType        string    `json:"flow_type"`        // topup / income / expense
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

// type PaginationMeta struct {
// 	CurrentPage  int `json:"current_page"`
// 	TotalPage    int `json:"total_page"`
// 	TotalRecords int `json:"total_records"`
// 	Limit        int `json:"limit"`
// }

type TransactionHistoryResponse struct {
	Transactions []TransactionHistoryItem `json:"transactions"`
	Meta         PaginationMeta           `json:"meta"`
}

type TransactionReportFilterParam struct {
	Type      string `form:"type" binding:"omitempty,oneof=income expense both"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type TransactionReportItem struct {
	Date         string `json:"date"`
	TotalIncome  int    `json:"total_income"`
	TotalExpense int    `json:"total_expense"`
}
