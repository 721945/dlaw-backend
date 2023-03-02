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
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	logger   *libs.Logger
}

func (u userService) WithTrx(trxHandle *gorm.DB) UserService {
	if trxHandle == nil {
		return u
	}

	u.userRepo = u.userRepo.WithTrx(trxHandle)
	return u
}

func (u userService) CreateUser(user models.User) (models.User, error) {
	return u.userRepo.CreateUser(user)
}

func (u userService) UpdateUser(id uint, user models.User) error {
	return u.userRepo.UpdateUser(id, user)
}

func (u userService) DeleteUser(id uint) error {
	return u.userRepo.DeleteUser(id)
}

func (u userService) GetUser(id uint) (*models.User, error) {
	return u.userRepo.GetUser(id)
}

func (u userService) GetUsers() ([]models.User, error) {
	return u.userRepo.GetUsers()
}

func (u userService) GetUserByEmail(email string) (user *models.User, err error) {
	return u.userRepo.GetUserByEmail(email)
}

func NewUserService(r repositories.UserRepository, logger *libs.Logger) UserService {
	return &userService{userRepo: r, logger: logger}
}
