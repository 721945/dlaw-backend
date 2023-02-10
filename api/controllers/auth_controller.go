package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	logger      libs.Logger
	userService services.UserService
	jwtService  services.JWTAuthService
}

func NewAuthController(logger libs.Logger) AuthController {
	return AuthController{logger: logger}
}

func (a AuthController) Login(c *gin.Context) {
	var input dtos.LoginDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.userService.GetUserByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	token := a.jwtService.GenerateToken(user)
	
	c.JSON(http.StatusOK, gin.H{"data": token})

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
