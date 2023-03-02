package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/utils"
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
	logger     *libs.Logger
	jwtService services.JWTAuthService
}

func NewUserController(service services.UserService, logger *libs.Logger, jwtService services.JWTAuthService) UserController {
	return UserController{service: service, logger: logger, jwtService: jwtService}
}

func (u UserController) GetUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)

	if err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := (u.service).GetUser(uint(id))

	if err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func (u UserController) GetMe(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization Header"})
		return
	}

	bearerToken := strings.Split(authHeader, " ")

	if len(bearerToken) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization Header"})
		return
	}

	token := bearerToken[1]

	id, err := u.jwtService.GetUserIDFromToken(token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := (u.service).GetUser(id)

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

	hashPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:     input.Email,
		Password:  hashPassword,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	//user, err =(u.service).CreateUser(user)
	user, err = (u.service).WithTrx(trxHandle).CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (u UserController) GetUsers(c *gin.Context) {
	users, err := (u.service).GetUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//dto := dtos.UserDto{}
	// Change from model to dto

	usersDto := make([]dtos.UserDto, len(users))

	for i, user := range users {
		usersDto[i] = dtos.UserDto{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.Firstname,
			LastName:  user.Lastname,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": usersDto})
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
	us, existed := c.Get("id")

	if existed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id in jwt token"})
		return
	}

	id := us.(uint)

	user := models.User{
		Email:     input.Email,
		Password:  input.Password,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	err := (u.service).WithTrx(trxHandle).UpdateUser(id, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
