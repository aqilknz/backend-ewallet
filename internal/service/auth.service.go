package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidInput       = errors.New("data input tidak valid")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar di sistem")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInternalServer     = errors.New("terjadi kesalahan pada server")
)

type AuthService struct {
	db       *pgxpool.Pool
	authRepo repository.AuthRepository
	rdb      *redis.Client
	mailer   pkg.Mailer
}

func NewAuthService(db *pgxpool.Pool, repo repository.AuthRepository, rdb *redis.Client, mailer pkg.Mailer) *AuthService {
	return &AuthService{
		db:       db,
		authRepo: repo,
		rdb:      rdb,
		mailer:   mailer,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dto.RegisterRequest) (dto.RegisterDataResponse, error) {
	var resData dto.RegisterDataResponse

	if !pkg.IsValidEmail(req.Email) {
		return resData, fmt.Errorf("%w: format email salah", ErrInvalidInput)
	}

	if len(req.Password) < 8 {
		return resData, fmt.Errorf("%w: password minimal 8 karakter", ErrInvalidInput)
	}

	exists, err := s.authRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return resData, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	if exists {
		return resData, ErrEmailAlreadyExists
	}

	hashedPassword, err := pkg.HashData(req.Password)
	if err != nil {
		return resData, fmt.Errorf("%w: gagal memproses password", ErrInternalServer)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return resData, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	defer tx.Rollback(ctx)

	newUser, err := s.authRepo.CreateUser(ctx, tx, req.Email, hashedPassword)
	if err != nil {
		return resData, fmt.Errorf("%w: gagal membuat user: %v", ErrInternalServer, err)
	}

	if err := s.authRepo.CreateProfile(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("%w: gagal membuat profile: %v", ErrInternalServer, err)
	}

	if err := s.authRepo.CreateWallet(ctx, tx, newUser.ID); err != nil {
		return resData, fmt.Errorf("%w: gagal membuat wallet: %v", ErrInternalServer, err)
	}

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

func (s *AuthService) LoginUser(ctx context.Context, req dto.LoginRequest) (string, bool, error) {
	if req.Email == "" || req.Password == "" {
		return "", false, fmt.Errorf("%w: email dan password wajib diisi", ErrInvalidInput)
	}

	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" || errors.Is(err, pgx.ErrNoRows) {
			return "", false, ErrInvalidCredentials
		}
		return "", false, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	match, err := pkg.VerifyHash(req.Password, user.Password)
	if err != nil || !match {
		return "", false, ErrInvalidCredentials
	}

	token, err := pkg.GenerateToken(int(user.ID))
	if err != nil {
		return "", false, fmt.Errorf("%w: gagal membuat sesi login", ErrInternalServer)
	}
	hasPin := user.Pin != ""

	return token, hasPin, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int, token string) error {
	if token == "" {
		return fmt.Errorf("%w: token kosong", ErrInvalidInput)
	}
	expiresIn := 12 * time.Hour

	err := s.authRepo.AddTokenToBlacklist(ctx, userID, token, expiresIn)
	if err != nil {
		return fmt.Errorf("%w: gagal memproses logout", ErrInternalServer)
	}

	return nil
}

func (s *AuthService) CreatePin(ctx context.Context, userID int, req dto.CreatePinRequest) error {
	currentPin, err := s.authRepo.GetUserPin(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user tidak ditemukan")
		}
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	if currentPin != "" {
		return fmt.Errorf("%w: PIN sudah dibuat, silakan gunakan menu ubah PIN", ErrInvalidInput)
	}

	hashedPin, err := pkg.HashData(req.Pin)
	if err != nil {
		return fmt.Errorf("%w: gagal memproses PIN", ErrInternalServer)
	}

	if err := s.authRepo.CreatePin(ctx, userID, hashedPin); err != nil {
		return fmt.Errorf("%w: gagal menyimpan PIN", ErrInternalServer)
	}

	return nil
}

func (s *AuthService) CheckEmail(ctx context.Context, req dto.CheckEmailRequest) error {
	exists, err := s.authRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	if !exists {
		return errors.New("email tidak ditemukan di sistem")
	}
	return nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest) error {
	hashedPassword, err := pkg.HashData(req.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: gagal mengamankan password baru", ErrInternalServer)
	}

	err = s.authRepo.UpdatePassword(ctx, req.Email, hashedPassword)
	if err != nil {
		if err.Error() == "gagal menyimpan perubahan ke database (email tidak ditemukan)" {
			return err
		}
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	exists, err := s.authRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	if !exists {
		return errors.New("email tidak terdaftar di sistem")
	}

	cooldownKey := fmt.Sprintf("cooldown_otp:%s", req.Email)
	if s.rdb.Get(ctx, cooldownKey).Err() == nil {
		return errors.New("harap tunggu 60 detik sebelum meminta OTP baru")
	}

	otpCode, err := pkg.GenerateOTP()
	if err != nil {
		return fmt.Errorf("gagal membuat OTP: %w", err)
	}

	cacheKey := fmt.Sprintf("otp:%s", req.Email)
	err = s.rdb.Set(ctx, cacheKey, otpCode, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("gagal menyimpan OTP ke cache: %w", err)
	}

	s.rdb.Set(ctx, cooldownKey, "true", 60*time.Second)

	err = s.mailer.SendResetPasswordOTP(req.Email, otpCode)
	if err != nil {
		s.rdb.Del(ctx, cooldownKey)
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) error {
	cacheKey := fmt.Sprintf("otp:%s", req.Email)
	savedOTP, err := s.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("kode OTP sudah kedaluwarsa atau tidak valid")
	} else if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	// validasi otp
	if savedOTP != req.OTP {
		return errors.New("kode OTP yang Anda masukkan salah")
	}
	s.rdb.Del(ctx, cacheKey)
	allowKey := fmt.Sprintf("reset_allowed:%s", req.Email)
	s.rdb.Set(ctx, allowKey, "true", 5*time.Minute)

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	allowKey := fmt.Sprintf("reset_allowed:%s", req.Email)

	allowed, err := s.rdb.Get(ctx, allowKey).Result()
	if errors.Is(err, redis.Nil) || allowed != "true" {
		return errors.New("akses ditolak, silakan lakukan verifikasi OTP terlebih dahulu")
	}

	hashedPassword, err := pkg.HashData(req.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: gagal mengamankan password baru", ErrInternalServer)
	}
	err = s.authRepo.UpdatePassword(ctx, req.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	s.rdb.Del(ctx, allowKey)

	return nil
}
