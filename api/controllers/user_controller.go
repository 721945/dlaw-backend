package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/constants"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type UserController struct {
	service    services.UserService
	logger     libs.Logger
	jwtService services.JWTAuthService
}

func NewUserController(s services.UserService, logger libs.Logger) UserController {
	return UserController{service: s, logger: logger}
}

func (u UserController) GetUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)

	if err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.GetUser(uint(id))

	if err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func (u UserController) CreateUser(c *gin.Context) {
	var input dtos.CreateUserDto
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&input); err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Email:     input.Email,
		Password:  input.Password,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	user, err := u.service.WithTrx(trxHandle).CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (u UserController) UpdateUser(c *gin.Context) {
	var input dtos.UpdateUserDto
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&input); err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get id from jwt token
	tokenPrefix := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenPrefix, "Bearer ")
	id, err := u.jwtService.GetUserIDFromToken(token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//claims := u.jwtService.VerifyToken()
	user := models.User{
		Email:     input.Email,
		Password:  input.Password,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	err = u.service.WithTrx(trxHandle).UpdateUser(id, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
