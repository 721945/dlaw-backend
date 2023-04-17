package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type PermissionRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewPermissionRepository(logger *libs.Logger, db libs.Database) PermissionRepository {
	return PermissionRepository{logger: logger, db: db}
}

func (r *PermissionRepository) GetPermissions() (permissions []models.Permission, err error) {
	return permissions, r.db.DB.Find(&permissions).Error

}

func (r *PermissionRepository) GetPermission(id uuid.UUID) (permission *models.Permission, err error) {
	return permission, r.db.DB.First(&permission, id).Error
}

func (r *PermissionRepository) GetPermissionByName(name string) (permission *models.Permission, err error) {
	return permission, r.db.DB.Where("name = ?", name).First(&permission).Error
}

func (r *PermissionRepository) CreatePermission(permission models.Permission) (models.Permission, error) {
	return permission, r.db.DB.Create(&permission).Error
}

func (r *PermissionRepository) UpdatePermission(id uuid.UUID, permission models.Permission) error {
	return r.db.DB.Model(&models.Permission{}).Where("id = ?", id).Updates(permission).Error
}

func (r *PermissionRepository) DeletePermission(id uuid.UUID) error {
	return r.db.DB.Delete(&models.Permission{}, id).Error
}
