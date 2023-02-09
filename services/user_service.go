package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"gorm.io/gorm"
)

type UserService interface {
	WithTrx(trxHandle *gorm.DB) UserService
	CreateUser(user models.User) error
	UpdateUser() error
	DeleteUser(id string) error
	GetUser(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
}

type userService struct {
	repository repositories.UserRepository
	logger     libs.Logger
}

func (u userService) WithTrx(trxHandle *gorm.DB) UserService {
	if trxHandle == nil {
		return u
	}

	u.repository = u.repository.WithTrx(trxHandle)
	return u
}

func (u userService) CreateUser(user models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) UpdateUser() error {
	//TODO implement me
	panic("implement me")
}

func (u userService) DeleteUser(id string) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) GetUser(id uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userService) GetUsers() ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService(r repositories.UserRepository, logger libs.Logger) UserService {
	return &userService{repository: r, logger: logger}
}
