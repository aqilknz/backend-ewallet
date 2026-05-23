package controller

import (
	"net/http"
	"strings"

	"github.com/aqilknz/backend-ewallet/internal/dto"
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
//	@Summary		Register a user
//	@Description	Create a new user account for E-Wallet
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body body dto.RegisterRequest true "register payload"
//	@Success		201	{object}	dto.AuthResponse
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	userData, err := ac.authService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "tidak valid") || strings.Contains(msg, "sudah terdaftar") {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Terjadi kegagalan server internal",
			"detail": msg,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		Message: "Registrasi berhasil, akun telah dibuat",
		Data:    userData,
	})
}

// User Login
//
//	@Summary		Login a user
//	@Description	Authenticate user and get JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body body dto.LoginRequest true "login payload"
//	@Success		200	{object}	dto.AuthResponse
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/auth [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi"})
		return
	}

	token, err := ac.authService.LoginUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Message: "Login berhasil",
		Token:   token,
	})
}
