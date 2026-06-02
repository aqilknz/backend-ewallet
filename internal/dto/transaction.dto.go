package dto

import "time"

type CheckPinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}

type TopUpRequest struct {
	Amount          int `json:"amount" binding:"required,min=10000"`
	PaymentMethodID int `json:"payment_method_id" binding:"required"`
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

type TransferRequest struct {
	Receiver_ID string `json:"receiver_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=10000"`
	Notes       string `json:"notes" binding:"omitempty,max=100"`
	Pin         string `json:"pin" binding:"required,len=6"`
}

type TransferResponse struct {
	TransactionID int       `json:"transaction_id"`
	SenderID      int       `json:"sender_id"`
	ReceiverID    int       `json:"receiver_id"`
	Amount        int       `json:"amount"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionHistoryFilterParam struct {
	Search string `form:"search" binding:"omitempty"`
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type TransactionHistoryItem struct {
	ID              int       `json:"id"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	FlowType        string    `json:"flow_type"`
	Description     string    `json:"description"`
	Phone           string    `json:"phone,omitempty"`
	Photo           string    `json:"photo,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransactionHistoryResponse struct {
	Transactions []TransactionHistoryItem `json:"transactions"`
	Meta         PaginationMeta           `json:"meta"`
}

type TransactionReportFilterParam struct {
	Type      string `form:"type" binding:"omitempty,oneof=income expense both"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

type TransactionReportItem struct {
	Date         string `json:"date"`
	TotalIncome  int    `json:"total_income"`
	TotalExpense int    `json:"total_expense"`
}
