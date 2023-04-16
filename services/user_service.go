package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	WithTrx(trxHandle *gorm.DB) UserService
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id uuid.UUID, user models.User) error
	DeleteUser(id uuid.UUID) error
	GetUser(id uuid.UUID) (*models.User, error)
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

func (u userService) UpdateUser(id uuid.UUID, user models.User) error {
	return u.userRepo.UpdateUser(id, user)
}

func (u userService) DeleteUser(id uuid.UUID) error {
	return u.userRepo.DeleteUser(id)
}

func (u userService) GetUser(id uuid.UUID) (*models.User, error) {
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
