package controller

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/response"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// Get User Profile
//
//	@Summary        Get user profile
//	@Description    Detailed profile information of the login user
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Success        200 {object}    dto.Response[dto.UserProfileResponse]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        404 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /users/profile [get]
func (uc *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)

	profile, err := uc.userService.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.JSONNotFound(ctx, "Profil tidak ditemukan", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, profile, "Profil berhasil diambil")
}

// Get User Dashboard
//
//	@Summary        Get user dashboard
//	@Description    Detailed dashboard information of user
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Success        200 {object}    dto.Response[dto.DashboardResponse]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        404 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /users/dashboard [get]
func (uc *UserController) GetDashboard(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)

	data, err := uc.userService.GetDashboard(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.JSONNotFound(ctx, "Data tidak ditemukan", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, data, "Data dashboard berhasil diambil")
}

// Edit Profile
//
//	@Summary        Edit user profile data
//	@Description    edit user detail fullname, phone and picture
//	@Tags           users
//	@Accept         mpfd
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          fullname    formData    string  false   "Update Fullname"
//	@Param          phone       formData    string  false   "Update Phone"
//	@Param          picture     formData    file    false   "Update Profile Picture (Max 2MB)"
//	@Success        200         {object}    dto.Response[dto.UserProfileResponse]
//	@Failure        400         {object}    dto.Response[any]
//	@Failure        401         {object}    dto.Response[any]
//	@Failure        422         {object}    dto.Response[any]
//	@Failure        500         {object}    dto.Response[any]
//	@Router         /users/profile [patch]
func (uc *UserController) EditProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)

	var req dto.EditProfileRequest

	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		response.JSONBadRequest(ctx, "Data input form tidak valid")
		return
	}

	var pictureURL *string

	if req.Picture != nil {
		const maxUploadSize = 2 * 1024 * 1024
		if req.Picture.Size > maxUploadSize {
			response.JSONUnprocessableEntity(ctx, "Ukuran file terlalu besar", "Maksimal 2MB")
			return
		}

		ext := strings.ToLower(filepath.Ext(req.Picture.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			response.JSONUnprocessableEntity(ctx, "Format file tidak didukung", "Gunakan .jpg, .jpeg, atau .png")
			return
		}

		filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
		dst := filepath.Join("public", "img", "profiles", filename)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			response.JSONInternalServerError(ctx, "Gagal membuat direktori penyimpanan gambar")
			return
		}

		if err := ctx.SaveUploadedFile(req.Picture, dst); err != nil {
			response.JSONInternalServerError(ctx, "Gagal menyimpan gambar profil")
			return
		}

		generatedURL := "/img/profiles/" + filename
		pictureURL = &generatedURL
	}

	data, err := uc.userService.EditProfile(ctx.Request.Context(), userID, req, pictureURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnprocessableEntity(ctx, "Gagal memperbarui profil", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, data, "Profil berhasil diperbarui")
}

// Update Password
//
//	@Summary        Change user password
//	@Description    Update the account password by verifying the old password first
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          body body       dto.EditPasswordRequest true "Update password payload"
//	@Success        200 {object}    dto.Response[any]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        422 {object}    dto.Response[any]
//	@Router         /users/profile/password [patch]
func (uc *UserController) EditPassword(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.EditPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Format input password tidak valid")
		return
	}
	if err := uc.userService.EditPassword(ctx.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidInput) || errors.Is(err, service.ErrInvalidCredentials) {
			response.JSONUnprocessableEntity(ctx, "Gagal mengubah password", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess[any](ctx, nil, "Password berhasil diubah")
}

// Update PIN
//
//	@Summary        Setup or change user transaction PIN
//	@Description    Update the 6-digit PIN used for validating transactions
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          body body       dto.EditPinRequest true "Update PIN payload"
//	@Success        200 {object}    dto.Response[any]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        422 {object}    dto.Response[any]
//	@Router         /users/profile/pin [patch]
func (uc *UserController) EditPin(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.EditPinRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Format input PIN tidak valid")
		return
	}

	if err := uc.userService.EditPin(ctx.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnprocessableEntity(ctx, "Gagal mengubah PIN", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess[any](ctx, nil, "Pin berhasil diubah")
}

// Check PIN
//
//	@Summary        Verify PIN for transaction
//	@Description    Check if the provided 6-digit PIN matches the user's PIN before authorizing a transaction
//	@Tags           transaction
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          body body       dto.CheckPinRequest true "Check PIN payload"
//	@Success        200 {object}    dto.Response[any]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /users/transaction/checkpin [post]
func (uc *UserController) CheckPin(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.CheckPinRequest

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.JSONBadRequest(ctx, "Format input PIN tidak valid")
		return
	}
	if err := uc.userService.CheckPin(ctx.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidInput) || errors.Is(err, service.ErrInvalidCredentials) {
			response.JSONUnauthorized(ctx, "Akses ditolak", "PIN yang Anda masukkan salah")
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess[any](ctx, nil, "PIN valid")
}

// Find Receivers
//
//	@Summary        Find receivers for transfer
//	@Description    Search other users by name, email, or phone with pagination
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          search  query   string  false   "Search by name, email, phone"
//	@Param          page    query   int     false   "Page number"       default(1)
//	@Param          limit   query   int     false   "Items per page"    default(10)
//	@Success        200     {object}    dto.Response[dto.ReceiverListResponse]
//	@Failure        400     {object}    dto.Response[any]
//	@Failure        500     {object}    dto.Response[any]
//	@Router         /users/receivers [get]
func (uc *UserController) FindReceivers(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var param dto.ReceiverFilterParam

	if err := ctx.ShouldBindQuery(&param); err != nil {
		response.JSONBadRequest(ctx, "Parameter query pencarian tidak valid")
		return
	}

	result, err := uc.userService.FindReceivers(ctx.Request.Context(), userID, param)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, result, "Berhasil mengambil data penerima")
}
