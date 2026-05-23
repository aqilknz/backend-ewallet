package controller

import (
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
//	@Summary		Get user profile
//	@Description	Detailed profile information of the login user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Failure		500	{object}	dto.Response
//	@Router			/users/profile [get]
func (uc *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)

	profile, err := uc.userService.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, profile, "Profil berhasil diambil")
}

// Get User Dashboard
//
//	@Summary		Get user dashboard
//	@Description	Detailed dashboard information of user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Failure		500	{object}	dto.Response
//	@Router			/users/dashboard [get]
func (uc *UserController) GetDashboard(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)

	data, err := uc.userService.GetDashboard(ctx.Request.Context(), userID)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, data, "Data dashboard berhasil diambil")
}

// Update Profile
//
//	@Summary		Update user profile data
//	@Description	Modify user details such as name, phone number, etc.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body dto.EditProfileRequest true "Update profile payload"
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Router			/users/profile [put]
func (uc *UserController) EditProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.EditProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, err.Error())
		return
	}

	data, err := uc.userService.EditProfile(ctx.Request.Context(), userID, req)
	if err != nil {
		response.JSONInternalServerError(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, data, "Profil berhasil diperbarui")
}

// Update Password
//
//	@Summary		Change user password
//	@Description	Update the account password by verifying the old password first
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body dto.EditPasswordRequest true "Update password payload"
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Router			/users/profile/password [patch]
func (uc *UserController) EditPassword(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.EditPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, err.Error())
		return
	}

	if err := uc.userService.EditPassword(ctx.Request.Context(), userID, req); err != nil {
		response.JSONBadRequest(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, nil, "Password berhasil diubah")
}

// Update PIN
//
//	@Summary		Setup or change user transaction PIN
//	@Description	Update the 6-digit PIN used for validating transactions
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body dto.EditPinRequest true "Update PIN payload"
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Router			/users/profile/pin [patch]
func (uc *UserController) EditPin(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.EditPinRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JSONBadRequest(ctx, "PIN harus 6 digit angka")
		return
	}

	if err := uc.userService.EditPin(ctx.Request.Context(), userID, req); err != nil {
		response.JSONBadRequest(ctx, err.Error())
		return
	}

	response.JSONSuccess(ctx, nil, "Pin berhasil diubah")
}

// Check PIN
//
//	@Summary		Verify PIN for transaction
//	@Description	Check if the provided 6-digit PIN matches the user's PIN before authorizing a transaction
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body body dto.CheckPinRequest true "Check PIN payload"
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	dto.Response
//	@Failure		401	{object}	dto.Response
//	@Router			/users/transaction/checkpin [post]
func (uc *UserController) CheckPin(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int)
	var req dto.CheckPinRequest

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.JSONBadRequest(ctx, err.Error())
		return
	}
	if err := uc.userService.CheckPin(ctx.Request.Context(), userID, req); err != nil {
		response.JSONUnauthorized(ctx, "Akses ditolak", err.Error())
		return
	}

	response.JSONSuccess(ctx, nil, "PIN valid")
}
