package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type CreatePinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}

type RegisterDataResponse struct {
	ID int `json:"id"`

	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
	HasPin  bool   `json:"has_pin"`
}
