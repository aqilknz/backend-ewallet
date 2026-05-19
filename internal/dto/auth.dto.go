package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterDataResponse struct {
	ID int `json:"id"`

	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Data    any    `json:"data,omitempty"`
}
