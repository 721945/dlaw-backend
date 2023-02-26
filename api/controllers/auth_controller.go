package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/util"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
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

	if !util.CheckPasswordHash(input.Password, user.Password) {
		_ = c.Error(libs.ErrUnauthorized)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	//
	token := ctrl.jwtService.GenerateToken(*user)

	c.JSON(http.StatusOK, gin.H{"data": token})

}
