package controller

import (
	"errors"
	"fmt"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/response"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// User Register
//
//	@Summary        Register a user
//	@Description    create a new user for e-wallet
//	@Tags           auth
//	@Accept         json
//	@Produce        json
//	@Param          body body dto.RegisterRequest true "register payload"
//	@Success        201 {object}    dto.Response[dto.RegisterDataResponse]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        409 {object}    dto.Response[any]
//	@Failure        422 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Format input tidak sesuai. Pastikan semua data terisi dengan benar.")
		return
	}

	userData, err := ac.authService.RegisterUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.JSONConflict(ctx, "Registrasi ditolak", err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnprocessableEntity(ctx, "Data tidak valid", err.Error())
			return
		}

		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONCreated(ctx, userData, "Registrasi berhasil, akun telah dibuat")
}

// User Login
//
//	@Summary        Login a user
//	@Description    Authenticate user and get JWT token dengan pengecekan status PIN
//	@Tags           auth
//	@Accept         json
//	@Produce        json
//	@Param          body body dto.LoginRequest true "login payload"
//	@Success        200 {object}    dto.Response[dto.AuthResponse]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /auth [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Email dan password wajib diisi dengan benar")
		return
	}

	token, hasPin, err := ac.authService.LoginUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnauthorized(ctx, "Gagal login", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}
	message := "Login berhasil"
	if !hasPin {
		message = "silahkan buat pin terlebih dahulu"
	}
	fmt.Println("DEBUG: Fungsi Login berhasil sampai bawah! Token:", token != "", "HasPin:", hasPin)

	response.JSONSuccess(ctx, dto.AuthResponse{
		Token:  token,
		HasPin: hasPin,
	}, message)
}

// User Logout
//
//	@Summary        Logout User
//	@Description    Invalidate current JWT token by adding it to Redis blacklist
//	@Tags           auth
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Success        200 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /auth/logout [delete]
func (ac *AuthController) Logout(ctx *gin.Context) {
	tokenString := ctx.GetString("token")
	if tokenString == "" {
		response.JSONUnauthorized(ctx, "Akses ditolak", "Token tidak ditemukan di sesi")
		return
	}

	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		response.JSONUnauthorized(ctx, "Sesi tidak valid", "Gagal mendapatkan ID pengguna")
		return
	}
	userID := userIDRaw.(int)

	err := ac.authService.Logout(ctx.Request.Context(), userID, tokenString)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess[any](ctx, nil, "Berhasil logout")
}

// Create PIN
//
//	@Summary        Create user transaction PIN
//	@Description    Setup the initial 6-digit PIN for a new user
//	@Tags           auth
//	@Accept         json
//	@Produce        json
//	@Security       ApiKeyAuth
//	@Param          body body       dto.CreatePinRequest true "Create PIN payload"
//	@Success        201 {object}    dto.Response[any]
//	@Failure        400 {object}    dto.Response[any]
//	@Failure        401 {object}    dto.Response[any]
//	@Failure        422 {object}    dto.Response[any]
//	@Failure        500 {object}    dto.Response[any]
//	@Router         /auth/create-pin [post]
func (ac *AuthController) CreatePin(ctx *gin.Context) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		response.JSONUnauthorized(ctx, "Akses ditolak", "Token tidak valid atau tidak ditemukan")
		return
	}
	userID := userIDRaw.(int)

	var req dto.CreatePinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Format PIN tidak valid. Pastikan PIN berupa 6 digit angka.")
		return
	}

	if err := ac.authService.CreatePin(ctx.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnprocessableEntity(ctx, "Pembuatan PIN gagal", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}
	response.JSONCreated[any](ctx, nil, "PIN berhasil dibuat")
}
