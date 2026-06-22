package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPinNotSet = errors.New("silahkan buat pin terlebih dahulu")
)

type TransactionService struct {
	txRepo repository.TransactionRepository
	rdb    *redis.Client
}

func NewTransactionService(txRepo repository.TransactionRepository, rdb *redis.Client) *TransactionService {
	return &TransactionService{txRepo: txRepo, rdb: rdb}
}

func (s *TransactionService) TopUp(ctx context.Context, userID int, req dto.TopUpRequest) (dto.TopUpResponse, error) {
	tax := 4000
	discount := 0

	subTotal := req.Amount + tax - discount

	res, err := s.txRepo.CreateTopUp(ctx, userID, req, tax, discount, subTotal)
	if err != nil {
		return dto.TopUpResponse{}, fmt.Errorf("gagal memproses top up: %v", err)
	}
	s.rdb.Del(ctx, fmt.Sprintf("dashboard:%d", userID))

	return res, nil
}

func (s *TransactionService) Transfer(ctx context.Context, senderID int, req dto.TransferRequest) (dto.TransferResponse, error) {
	hashedPin, err := s.txRepo.GetSenderPin(ctx, senderID)
	if err != nil {
		return dto.TransferResponse{}, fmt.Errorf("gagal memverifikasi pengguna: %v", err)
	}

	if hashedPin == "" {
		return dto.TransferResponse{}, ErrPinNotSet
	}

	match, err := pkg.VerifyHash(req.Pin, hashedPin)
	if err != nil || !match {
		return dto.TransferResponse{}, errors.New("PIN transaksi salah")
	}

	receiverID, err := strconv.Atoi(req.Receiver_ID)
	if err != nil {
		return dto.TransferResponse{}, errors.New("format ID penerima tidak valid")
	}

	if senderID == receiverID {
		return dto.TransferResponse{}, errors.New("tidak dapat melakukan transfer ke akun sendiri")
	}

	res, err := s.txRepo.CreateTransfer(ctx, senderID, receiverID, req.Amount, req.Notes)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	s.rdb.Del(ctx, fmt.Sprintf("dashboard:%d", senderID), fmt.Sprintf("dashboard:%d", receiverID))

	return res, nil
}

func (s *TransactionService) GetHistory(ctx context.Context, userID int, param dto.TransactionHistoryFilterParam) (dto.TransactionHistoryResponse, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	offset := (param.Page - 1) * param.Limit

	histories, totalRecords, err := s.txRepo.GetHistory(ctx, userID, param.Search, param.Limit, offset)
	if err != nil {
		return dto.TransactionHistoryResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
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

	report, err := s.txRepo.GetReport(ctx, userID, param)
	if err != nil {
		return nil, fmt.Errorf("%w: gagal mengambil laporan", ErrInternalServer)
	}

	return report, nil
}
