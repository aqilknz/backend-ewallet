package controller

import (
	"errors"
	"strings"

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
//		@Summary		Register a user
//		@Description	Create a new user account for E-Wallet
//		@Tags			auth
//		@Accept			json
//		@Produce		json
//		@Param			body body dto.RegisterRequest true "register payload"
//		@Success		201	{object}	dto.Response[dto.RegisterDataResponse]
//		@Failure		400	{object}	dto.Response[any]
//	 	@Failure		409 {object}	dto.Response[any]
//		@Failure		500	{object}	dto.Response[any]
//		@Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Data input tidak valid: "+err.Error())
		return
	}

	userData, err := ac.authService.RegisterUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {

			response.JSONConflict(ctx, "Registrasi ditolak", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONCreated(ctx, userData, "Registrasi berhasil, akun telah dibuat")
}

// User Login
//
//	@Summary		Login a user
//	@Description	Authenticate user and get JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body body dto.LoginRequest true "login payload"
//	@Success		200	{object}	dto.Response[dto.AuthResponse]
//	@Failure		400	{object}	dto.Response[any]
//	@Failure		401	{object}	dto.Response[any]
//	@Router			/auth [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "Email dan password wajib diisi")
		return
	}

	token, err := ac.authService.LoginUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrInvalidInput) {
			response.JSONUnauthorized(ctx, "Gagal login", err.Error())
			return
		}
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, dto.AuthResponse{Token: token}, "Login berhasil")
}

// User Logout
//
//	@Summary		Logout User
//	@Description	Invalidate current JWT token by adding it to database blacklist
//	@Tags			auth
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200 {object}	dto.Response[any]
//	@Failure		401 {object}	dto.Response[any]
//	@Failure		500 {object}	dto.Response[any]
//	@Router			/auth/logout [delete]
func (ac *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.JSONUnauthorized(ctx, "Akses ditolak", "Header Authorization tidak ditemukan")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	err := ac.authService.Logout(ctx.Request.Context(), tokenString)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess[any](ctx, nil, "Berhasil logout")
}
