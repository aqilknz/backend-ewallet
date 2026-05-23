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
		return resData, errors.New("format email tidak valid")
	}

	// Cek Duplikasi Email
	exists, err := s.authRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return resData, fmt.Errorf("error database: %v", err)
	}
	if exists {
		return resData, errors.New("email sudah terdaftar")
	}

	// cek panjang password
	if len(req.Password) < 6 {
		return resData, errors.New("password harus memiliki minimal 6 karakter")
	}

	// Hash Password
	hashedPassword, err := pkg.HashData(req.Password)
	if err != nil {
		return resData, errors.New("gagal memproses password")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return resData, err
	}
	defer tx.Rollback(ctx)

	// Simpan ke Tabel Users
	newUser, err := s.authRepo.CreateUser(ctx, tx, req.Email, hashedPassword)
	if err != nil {
		return resData, fmt.Errorf("gagal membuat user: %v", err)
	}

	// Simpan ke Tabel Profiles & Wallets dengan ID Baru
	if err := s.authRepo.CreateProfile(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("gagal membuat profil: %v", err)
	}

	if err := s.authRepo.CreateWallet(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("gagal membuat dompet: %v", err)
	}

	// Commit Transaksi jika semua sukses
	if err := tx.Commit(ctx); err != nil {
		return resData, err
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
		return "", errors.New("email dan password tidak boleh kosong")
	}

	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// Cek Validitas Password
	match, err := pkg.VerifyHash(req.Password, user.Password)
	if err != nil || !match {
		return "", errors.New("email atau password salah")
	}

	// Buat JWT Token
	token, err := pkg.GenerateToken(int(user.ID))
	if err != nil {
		return "", errors.New("gagal membuat sesi login")
	}

	return token, nil
}
