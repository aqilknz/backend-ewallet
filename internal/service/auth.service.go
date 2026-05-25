package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 1. Deklarasi Sentinel Errors
var (
	ErrInvalidInput       = errors.New("data input tidak valid")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar di sistem")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInternalServer     = errors.New("terjadi kesalahan pada server")
)

type AuthService struct {
	db       *pgxpool.Pool
	authRepo repository.AuthRepository
}

func NewAuthService(db *pgxpool.Pool, repo repository.AuthRepository) *AuthService {
	return &AuthService{db: db, authRepo: repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dto.RegisterRequest) (dto.RegisterDataResponse, error) {
	var resData dto.RegisterDataResponse

	// Validasi Format Email
	if !pkg.IsValidEmail(req.Email) {
		return resData, fmt.Errorf("%w: format email salah", ErrInvalidInput)
	}

	// cek panjang password
	if len(req.Password) < 6 {
		return resData, fmt.Errorf("%w: password minimal 6 karakter", ErrInvalidInput)
	}

	// Cek Duplikasi Email
	exists, err := s.authRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return resData, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	if exists {
		return resData, ErrEmailAlreadyExists // Kembalikan sentinel error secara langsung
	}

	// Hash Password
	hashedPassword, err := pkg.HashData(req.Password)
	if err != nil {
		return resData, fmt.Errorf("%w: gagal memproses password", ErrInternalServer)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return resData, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	defer tx.Rollback(ctx)

	// Simpan ke Tabel Users
	newUser, err := s.authRepo.CreateUser(ctx, tx, req.Email, hashedPassword)
	if err != nil {
		return resData, fmt.Errorf("%w: gagal membuat user: %v", ErrInternalServer, err)
	}

	// Simpan ke Tabel Profiles dan Wallets dengan ID Baru
	if err := s.authRepo.CreateProfile(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("%w: gagal membuat profil: %v", ErrInternalServer, err)
	}

	if err := s.authRepo.CreateWallet(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("%w: gagal membuat dompet: %v", ErrInternalServer, err)
	}

	// Commit Transaksi jika semua sukses
	if err := tx.Commit(ctx); err != nil {
		return resData, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	resData = dto.RegisterDataResponse{
		ID:        int(newUser.ID),
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return resData, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req dto.LoginRequest) (string, error) {
	// validasi input kosong
	if req.Email == "" || req.Password == "" {
		return "", fmt.Errorf("%w: email dan password wajib diisi", ErrInvalidInput)
	}

	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Cek Validitas Password
	match, err := pkg.VerifyHash(req.Password, user.Password)
	if err != nil || !match {
		return "", ErrInvalidCredentials
	}

	// buat JWT Token
	token, err := pkg.GenerateToken(int(user.ID))
	if err != nil {
		return "", fmt.Errorf("%w: gagal membuat sesi login", ErrInternalServer)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("%w: token kosong", ErrInvalidInput)
	}

	err := s.authRepo.AddTokenToBlacklist(ctx, token)
	if err != nil {
		return fmt.Errorf("%w: gagal memproses logout", ErrInternalServer)
	}

	return nil
}
