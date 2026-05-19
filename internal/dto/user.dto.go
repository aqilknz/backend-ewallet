package dto

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
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}

type EditPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type EditPinRequest struct {
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin" binding:"required,len=6"`
}
