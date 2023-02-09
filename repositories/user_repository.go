package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTrx(trxHandle *gorm.DB) UserRepository
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id uint, user models.User) error
	DeleteUser(id uint) error
	GetUser(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
}

type userRepository struct {
	db     libs.Database
	logger libs.Logger
}

func NewUserRepository(db *libs.Database, logger libs.Logger) UserRepository {
	return &userRepository{db: *db, logger: logger}
}

func (r *userRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		return r
	}

	r.db.DB = trxHandle
	return r
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) UpdateUser(id uint, user models.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) DeleteUser(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetUser(id uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetUsers() ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}
