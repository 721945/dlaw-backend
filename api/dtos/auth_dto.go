package dtos

type LoginDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ForgetPasswordDto struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordDto struct {
	Email string `json:"email" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

type ChangePasswordDto struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}
