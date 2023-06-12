package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

	id, err := uuid.Parse(paramID)

	if err != nil {
		//u.logger.Error(err)
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

func (u UserController) GetMe(c *gin.Context) {
	id, isExisted := c.Get("id")

	if !isExisted {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	user, err := (u.service).GetUser(id.(uuid.UUID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (u UserController) CreateUser(c *gin.Context) {
	var input dtos.CreateUserDto
	//trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organization := c.GetHeader("x-organization")

	if organization == "" {
		_ = c.Error(libs.ErrBadRequest)
		return
	}

	hashPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        input.Email,
		Password:     hashPassword,
		Firstname:    input.FirstName,
		Lastname:     input.LastName,
		Organization: organization,
	}

	user, err = (u.service).CreateUser(user)
	//user, err = (u.service).WithTrx(trxHandle).CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user.ID.String()})
}

func (u UserController) GetUsers(c *gin.Context) {
	organization := c.GetHeader("x-organization")

	if organization == "" {
		_ = c.Error(libs.ErrBadRequest)
		return
	}

	users, err := (u.service).GetUsers(organization)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//dto := dtos.UserDto{}
	// Change from model to dto

	usersDto := make([]dtos.UserDto, len(users))

	for i, user := range users {
		usersDto[i] = dtos.UserDto{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.Firstname,
			LastName:  user.Lastname,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": usersDto})
}

func (u UserController) UpdateUser(c *gin.Context) {
	var input dtos.UpdateUserDto
	//trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&input); err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get id from jwt token
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:     input.Email,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	err = (u.service).UpdateUser(id, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// UpdateCurrentUser

func (u UserController) UpdateCurrentUser(c *gin.Context) {
	var input dtos.UpdateUserDto
	if err := c.ShouldBindJSON(&input); err != nil {
		//u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get id from jwt token
	authId, existed := c.Get("id")

	if !existed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id in jwt token"})
		return
	}

	id := authId.(uuid.UUID)

	user := models.User{
		Email:     input.Email,
		Firstname: input.FirstName,
		Lastname:  input.LastName,
	}

	err := (u.service).UpdateUser(id, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
