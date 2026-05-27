package dto

import "mime/multipart"

type UserProfileResponse struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}

type DashboardResponse struct {
	Balance int `json:"balance"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

type EditProfileRequest struct {
	Fullname *string               `form:"fullname"`
	Phone    *string               `form:"phone"`
	Picture  *multipart.FileHeader `form:"picture" binding:"omitempty"`
}

type EditPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type EditPinRequest struct {
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin" binding:"required,len=6"`
}

type ReceiverFilterParam struct {
	Search string `form:"search"`
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type ReceiverResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPage    int `json:"total_page"`
	TotalRecords int `json:"total_records"`
	Limit        int `json:"limit"`
}

type ReceiverListResponse struct {
	Receivers []ReceiverResponse `json:"receivers"`
	Meta      PaginationMeta     `json:"meta"`
}
