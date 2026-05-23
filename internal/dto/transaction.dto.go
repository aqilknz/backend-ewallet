package dto

type CheckPinRequest struct {
	Pin string `json: pin binding:"required, len=6, numeric"`
}
