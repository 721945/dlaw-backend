package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"gorm.io/gorm"
)

type UserService interface {
	WithTrx(trxHandle *gorm.DB) UserService
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id uint, user models.User) error
	DeleteUser(id uint) error
	GetUser(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
	GetUserByEmail(email string) (models.User, error)
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

func (u userService) CreateUser(user models.User) (models.User, error) {
	return u.repository.CreateUser(user)
}

func (u userService) UpdateUser(id uint, user models.User) error {
	return u.repository.UpdateUser(id, user)
}

func (u userService) DeleteUser(id uint) error {
	return u.repository.DeleteUser(id)
}

func (u userService) GetUser(id uint) (*models.User, error) {
	return u.repository.GetUser(id)
}

func (u userService) GetUsers() ([]models.User, error) {
	return u.repository.GetUsers()
}

func (u userService) GetUserByEmail(email string) (models.User, error) {
	return u.repository.GetUserByEmail(email)
}

func NewUserService(r repositories.UserRepository, logger libs.Logger) UserService {
	return &userService{repository: r, logger: logger}
}
