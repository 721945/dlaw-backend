package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTrx(trxHandle *gorm.DB) UserRepository
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id uuid.UUID, user models.User) error
	DeleteUser(id uuid.UUID) error
	GetUser(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByEmailAndOrganization(email, organization string) (user *models.User, err error)
	GetUsers() ([]models.User, error)
	GetUsersByOrganization(organization string) ([]models.User, error)
}

type userRepository struct {
	db     libs.Database
	logger *libs.Logger
}

func NewUserRepository(db libs.Database, logger *libs.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		return r
	}

	r.db.DB = trxHandle
	return r
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	return user, r.db.DB.Create(&user).Error
}

func (r *userRepository) UpdateUser(id uuid.UUID, user models.User) error {
	return r.db.DB.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	return r.db.DB.Delete(&models.User{}, id).Error
}

func (r *userRepository) GetUser(id uuid.UUID) (user *models.User, err error) {
	return user, r.db.DB.First(&user, id).Error
}

func (r *userRepository) GetUsers() (users []models.User, err error) {
	return users, r.db.DB.Find(&users).Error
}
func (r *userRepository) GetUsersByOrganization(organization string) (users []models.User, err error) {
	return users, r.db.DB.Where("organization = ?", organization).Find(&users).Error
}

func (r *userRepository) GetUserByEmail(email string) (user *models.User, err error) {
	return user, r.db.DB.Where("email = ?", email).First(&user).Error
}

func (r *userRepository) GetUserByEmailAndOrganization(email, organization string) (user *models.User, err error) {
	return user, r.db.DB.Where("email = ? AND organization = ?", email, organization).First(&user).Error
}
