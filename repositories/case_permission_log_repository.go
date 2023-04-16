package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type CasePermissionLogRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewCasePermissionLogRepository(logger *libs.Logger, db libs.Database) CasePermissionLogRepository {
	return CasePermissionLogRepository{logger: logger, db: db}
}

func (r *CasePermissionLogRepository) GetPermissionLogs() (permissionLogs []models.CasePermissionLog, err error) {
	return permissionLogs, r.db.DB.Find(&permissionLogs).Error
}

func (r *CasePermissionLogRepository) GetPermissionLog(id uuid.UUID) (permissionLog *models.CasePermissionLog, err error) {
	return permissionLog, r.db.DB.First(&permissionLog, id).Error
}

func (r *CasePermissionLogRepository) CreatePermissionLog(permissionLog models.CasePermissionLog) (models.CasePermissionLog, error) {
	return permissionLog, r.db.DB.Create(&permissionLog).Error
}

func (r *CasePermissionLogRepository) UpdatePermissionLog(id uuid.UUID, permissionLog models.CasePermissionLog) error {
	return r.db.DB.Model(&models.CasePermissionLog{}).Where("id = ?", id).Updates(permissionLog).Error
}

func (r *CasePermissionLogRepository) DeletePermissionLog(id uuid.UUID) error {
	return r.db.DB.Delete(&models.CasePermissionLog{}, id).Error
}
