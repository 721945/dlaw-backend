package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AuthController struct {
	logger      *libs.Logger
	userService services.UserService
	jwtService  services.JWTAuthService
}

func NewAuthController(logger *libs.Logger, userService services.UserService, jwtService services.JWTAuthService) AuthController {
	return AuthController{logger: logger, userService: userService, jwtService: jwtService}
}

func (ctrl AuthController) Login(c *gin.Context) {
	var input dtos.LoginDto

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(libs.ErrBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.logger.Info(input.Email)

	user, err := ctrl.userService.GetUserByEmail(input.Email)

	if err != nil {
		_ = c.Error(libs.ErrUnauthorized)

		//panic(libs.ErrUnauthorized)
		//c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		_ = c.Error(libs.ErrUnauthorized)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	//
	token := ctrl.jwtService.GenerateToken(*user)

	c.JSON(http.StatusOK, gin.H{"data": token})

}

func (ctrl AuthController) ForgetPassword(c *gin.Context) {

	var input dtos.ForgetPasswordDto

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(libs.ErrBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userService.ForgetPassword(input.Email)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (ctrl AuthController) CheckOtp(c *gin.Context) {

	email := c.DefaultQuery("email", "")
	otp := c.DefaultQuery("otp", "")

	if email == "" || otp == "" {
		_ = c.Error(libs.ErrBadRequest)
		return
	}

	err := ctrl.userService.VerifyOTP(email, otp)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (ctrl AuthController) ResetPassword(c *gin.Context) {

	var input dtos.ResetPasswordDto

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(libs.ErrBadRequest)
		return
	}

	err := ctrl.userService.ResetPassword(input.Email, input.Otp, input.Password)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (ctrl AuthController) ChangePassword(c *gin.Context) {
	id, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	var input dtos.ChangePasswordDto

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(libs.ErrBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userService.ChangedPassword(id.(uuid.UUID), input.CurrentPassword, input.NewPassword)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
