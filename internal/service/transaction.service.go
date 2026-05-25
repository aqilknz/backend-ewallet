package service

import (
	"context"
	"errors"
	"math"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
)

type TransactionService struct {
	txRepo repository.TransactionRepository
}

func NewTransactionService(txRepo repository.TransactionRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo}
}

func (s *TransactionService) TopUp(ctx context.Context, userID int, req dto.TopUpRequest) (dto.TopUpResponse, error) {
	return s.txRepo.TopUp(ctx, userID, req)
}

func (s *TransactionService) Transfer(ctx context.Context, senderID int, req dto.TransferRequest) (dto.TransferResponse, error) {
	receiverID, err := s.txRepo.GetUserIDByEmail(ctx, req.ReceiverEmail)
	if err != nil {
		return dto.TransferResponse{}, errors.New("email penerima tidak ditemukan")
	}
	if senderID == receiverID {
		return dto.TransferResponse{}, errors.New("tidak bisa mentransfer ke akun sendiri")
	}

	return s.txRepo.Transfer(ctx, senderID, receiverID, req)
}

func (s *TransactionService) GetHistory(ctx context.Context, userID int, param dto.TransactionHistoryFilterParam) (dto.TransactionHistoryResponse, error) {
	offset := (param.Page - 1) * param.Limit

	histories, totalRecords, err := s.txRepo.GetHistory(ctx, userID, param.Search, param.Limit, offset)
	if err != nil {
		return dto.TransactionHistoryResponse{}, err
	}

	totalPage := int(math.Ceil(float64(totalRecords) / float64(param.Limit)))

	return dto.TransactionHistoryResponse{
		Transactions: histories,
		Meta: dto.PaginationMeta{
			CurrentPage:  param.Page,
			TotalPage:    totalPage,
			TotalRecords: totalRecords,
			Limit:        param.Limit,
		},
	}, nil
}

func (s *TransactionService) GetReport(ctx context.Context, userID int, param dto.TransactionReportFilterParam) ([]dto.TransactionReportItem, error) {
	if param.Type == "" || param.Type == "both" {
		param.Type = "all"
	}
	return s.txRepo.GetReport(ctx, userID, param)
}
