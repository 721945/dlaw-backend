package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
)

type PermissionLogRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewPermissionLogRepository(logger *libs.Logger, db libs.Database) PermissionLogRepository {
	return PermissionLogRepository{logger: logger, db: db}
}

func (r *PermissionLogRepository) GetPermissionLogs() (permissionLogs []models.CasePermissionLog, err error) {
	return permissionLogs, r.db.DB.Find(&permissionLogs).Error
}

func (r *PermissionLogRepository) GetPermissionLog(id uint) (permissionLog *models.CasePermissionLog, err error) {
	return permissionLog, r.db.DB.First(&permissionLog, id).Error
}

func (r *PermissionLogRepository) CreatePermissionLog(permissionLog models.CasePermissionLog) (models.CasePermissionLog, error) {
	return permissionLog, r.db.DB.Create(&permissionLog).Error
}

func (r *PermissionLogRepository) UpdatePermissionLog(id uint, permissionLog models.CasePermissionLog) error {
	return r.db.DB.Model(&models.CasePermissionLog{}).Where("id = ?", id).Updates(permissionLog).Error
}

func (r *PermissionLogRepository) DeletePermissionLog(id uint) error {
	return r.db.DB.Delete(&models.CasePermissionLog{}, id).Error
}
