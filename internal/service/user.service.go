package service

import (
	"context"
	"errors"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error) {
	return s.userRepo.GetProfile(ctx, userID)
}

func (s *UserService) GetDashboard(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	return s.userRepo.GetDashboard(ctx, userID)
}

func (s *UserService) EditProfile(ctx context.Context, userID int, req dto.EditProfileRequest) (dto.UserProfileResponse, error) {
	err := s.userRepo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}
	res := dto.UserProfileResponse{
		Email:    user.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}
	return res, nil
}

func (s *UserService) EditPassword(ctx context.Context, userID int, req dto.EditPasswordRequest) error {
	// Ambil data lama dari database
	oldHash, _, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Cocokkan password lama inputan user dengan hash di database
	match, _ := utils.CheckPassword(req.OldPassword, oldHash)
	if !match {
		return errors.New("password lama salah")
	}

	// Hash password baru
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	// Simpan ke database
	return s.userRepo.UpdatePassword(ctx, userID, newHash)
}

func (s *UserService) EditPin(ctx context.Context, userID int, req dto.EditPinRequest) error {
	_, oldPinHash, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Jika PIN sudah ada isinya, wajib cek kecocokan old_pin
	if oldPinHash != "" {
		match, _ := utils.CheckPassword(req.OldPin, oldPinHash)
		if !match {
			return errors.New("PIN lama salah")
		}
	}

	newPinHash, err := utils.HashPassword(req.NewPin)
	if err != nil {
		return errors.New("gagal memproses PIN baru")
	}

	return s.userRepo.UpdatePin(ctx, userID, newPinHash)
}
