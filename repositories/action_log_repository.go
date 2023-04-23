package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type ActionLogRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewActionLogRepository(logger *libs.Logger, db libs.Database) ActionLogRepository {
	return ActionLogRepository{logger: logger, db: db}
}

func (r *ActionLogRepository) GetActionLogs() (permissionLogs []models.ActionLog, err error) {
	return permissionLogs, r.db.DB.Find(&permissionLogs).Error
}

func (r *ActionLogRepository) GetActionLog(id uuid.UUID) (permissionLog *models.ActionLog, err error) {
	return permissionLog, r.db.DB.First(&permissionLog, id).Error
}

func (r *ActionLogRepository) CreateActionLog(permissionLog models.ActionLog) (models.ActionLog, error) {
	return permissionLog, r.db.DB.Create(&permissionLog).Error
}

func (r *ActionLogRepository) UpdateActionLog(id uuid.UUID, permissionLog models.ActionLog) error {
	return r.db.DB.Model(&models.ActionLog{}).Where("id = ?", id).Updates(permissionLog).Error
}

func (r *ActionLogRepository) DeleteActionLog(id uuid.UUID) error {
	return r.db.DB.Delete(&models.ActionLog{}, id).Error
}

// Get by folder id
func (r *ActionLogRepository) GetActionLogsByFolderId(folderId uuid.UUID) (permissionLogs []models.ActionLog, err error) {
	return permissionLogs, r.db.DB.Where("folder_id = ?", folderId).Find(&permissionLogs).Error
}
