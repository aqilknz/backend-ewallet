package service

import (
	"context"
	"errors"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg"
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
	if len(req.NewPassword) < 6 {
		return errors.New("password baru minimal harus 6 karakter")
	}
	oldHash, _, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	match, _ := pkg.VerifyHash(req.OldPassword, oldHash)
	if !match {
		return errors.New("password lama salah")
	}

	newHash, err := pkg.HashData(req.NewPassword)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	return s.userRepo.UpdatePassword(ctx, userID, newHash)
}

func (s *UserService) EditPin(ctx context.Context, userID int, req dto.EditPinRequest) error {
	if len(req.NewPin) < 6 {
		return errors.New("pin minimal 6 digit")
	}
	_, oldPinHash, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if oldPinHash != "" {
		match, _ := pkg.VerifyHash(req.OldPin, oldPinHash)
		if !match {
			return errors.New("PIN lama salah")
		}
	}

	newPinHash, err := pkg.HashData(req.NewPin)
	if err != nil {
		return errors.New("gagal memproses PIN baru")
	}

	return s.userRepo.UpdatePin(ctx, userID, newPinHash)
}

func (s *UserService) CheckPin(ctx context.Context, userID int, req dto.CheckPinRequest) error {
	_, PinHash, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	if PinHash == "" {
		return errors.New("Pin belum diatur, buatlah terlebih dahulu")
	}
	match, _ := pkg.VerifyHash(req.Pin, PinHash)
	if !match {
		return errors.New("Pin salah")
	}
	return nil
}
