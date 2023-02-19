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

func (a AuthController) Login(c *gin.Context) {
	var input dtos.LoginDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.logger.Info(input.Email)

	user, err := a.userService.GetUserByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !util.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	//
	token := a.jwtService.GenerateToken(*user)

	c.JSON(http.StatusOK, gin.H{"data": token})

}
