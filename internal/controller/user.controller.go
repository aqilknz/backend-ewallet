package controller

import (
	"net/http"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	profile, err := uc.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sukses", "data": profile})
}

func (uc *UserController) GetDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	data, err := uc.userService.GetDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sukses", "data": data})
}

func (uc *UserController) EditProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req dto.EditProfileRequest
	// var data dto.UserProfileResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}
	data, err := uc.userService.EditProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui", "data": data})
}

func (uc *UserController) EditPassword(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req dto.EditPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak lengkap atau tidak sesuai format"})
		return
	}

	if err := uc.userService.EditPassword(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

func (uc *UserController) EditPin(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req dto.EditPinRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN harus 6 digit angka"})
		return
	}

	if err := uc.userService.EditPin(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN berhasil diubah"})
}
