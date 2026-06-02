package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/repository"
	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound = errors.New("pengguna tidak ditemukan")
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error) {
	profile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserProfileResponse{}, ErrUserNotFound
		}
		return dto.UserProfileResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	return profile, nil
}

func (s *UserService) GetDashboard(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	data, err := s.userRepo.GetDashboard(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DashboardResponse{}, ErrUserNotFound
		}
		return dto.DashboardResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	return data, nil
}

func (s *UserService) EditProfile(ctx context.Context, userID int, req dto.EditProfileRequest, photoURL *string, deletePicture bool) (dto.UserProfileResponse, error) {
	if deletePicture || (photoURL != nil && *photoURL != "") {
		oldProfile, err := s.userRepo.GetProfile(ctx, userID)
		if err == nil && oldProfile.Photo != "" && !strings.HasPrefix(oldProfile.Photo, "http") {
			filename := filepath.Base(oldProfile.Photo)
			exactPath := filepath.Join("public", "img", "profiles", filename)
			_ = os.Remove(exactPath)
		}
	}

	if deletePicture && photoURL == nil {
		emptyPhoto := ""
		photoURL = &emptyPhoto
	}
	err := s.userRepo.EditProfile(ctx, userID, req.Fullname, req.Phone, photoURL)
	if err != nil {
		return dto.UserProfileResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	updatedProfile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		return dto.UserProfileResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	return updatedProfile, nil
}

func (s *UserService) EditPassword(ctx context.Context, userID int, req dto.EditPasswordRequest) error {
	if len(req.NewPassword) < 8 {
		return fmt.Errorf("%w: password baru minimal harus 8 karakter", ErrInvalidInput)
	}

	oldHash, _, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	match, _ := pkg.VerifyHash(req.OldPassword, oldHash)
	if !match {
		return fmt.Errorf("%w: password lama salah", ErrInvalidCredentials)
	}

	newHash, err := pkg.HashData(req.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: gagal memproses password baru", ErrInternalServer)
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		return fmt.Errorf("%w: gagal menyimpan password", ErrInternalServer)
	}
	return nil
}

func (s *UserService) EditPin(ctx context.Context, userID int, req dto.EditPinRequest) error {
	if len(req.NewPin) < 6 {
		return fmt.Errorf("%w: pin minimal 6 digit", ErrInvalidInput)
	}

	_, oldPinHash, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	if oldPinHash != "" {
		match, _ := pkg.VerifyHash(req.OldPin, oldPinHash)
		if !match {
			return fmt.Errorf("%w: PIN lama salah", ErrInvalidInput)
		}
	}

	newPinHash, err := pkg.HashData(req.NewPin)
	if err != nil {
		return fmt.Errorf("%w: gagal memproses PIN baru", ErrInternalServer)
	}

	if err := s.userRepo.UpdatePin(ctx, userID, newPinHash); err != nil {
		return fmt.Errorf("%w: gagal menyimpan PIN", ErrInternalServer)
	}
	return nil
}

func (s *UserService) CheckPin(ctx context.Context, userID int, req dto.CheckPinRequest) error {
	_, PinHash, err := s.userRepo.GetPasswordAndPin(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	if PinHash == "" {
		return fmt.Errorf("%w: PIN belum diatur, buatlah terlebih dahulu", ErrInvalidInput)
	}

	match, _ := pkg.VerifyHash(req.Pin, PinHash)
	if !match {
		return fmt.Errorf("%w: PIN salah", ErrInvalidCredentials)
	}

	return nil
}

func (s *UserService) FindReceivers(ctx context.Context, userID int, param dto.ReceiverFilterParam) (dto.ReceiverListResponse, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	offset := (param.Page - 1) * param.Limit

	receivers, totalRecords, err := s.userRepo.FindReceivers(ctx, userID, param.Search, param.Limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	totalPage := int(math.Ceil(float64(totalRecords) / float64(param.Limit)))

	return dto.ReceiverListResponse{
		Receivers: receivers,
		Meta: dto.PaginationMeta{
			CurrentPage:  param.Page,
			TotalPage:    totalPage,
			TotalRecords: totalRecords,
			Limit:        param.Limit,
		},
	}, nil
}
