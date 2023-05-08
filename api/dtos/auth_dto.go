package dtos

type LoginDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ForgetPasswordDto struct {
	Email string `json:"email" binding:"required"`
}

type OtpCheckDto struct {
	Email string `json:"email" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

type ResetPasswordDto struct {
	Email    string `json:"email" binding:"required"`
	Otp      string `json:"otp" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordDto struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}
