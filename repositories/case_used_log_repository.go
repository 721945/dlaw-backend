package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CaseUsedLogRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewCaseUsedLogRepository(logger *libs.Logger, db libs.Database) CaseUsedLogRepository {
	return CaseUsedLogRepository{logger: logger, db: db}
}

func (r CaseUsedLogRepository) GetCaseUsedLogs() (cases []models.CaseUsedLog, err error) {
	return cases, r.db.DB.Find(&cases).Error
}

func (r CaseUsedLogRepository) GetCaseUsedLog(id uuid.UUID) (mCase *models.CaseUsedLog, err error) {
	return mCase, r.db.DB.First(&mCase, id).Error
}

func (r CaseUsedLogRepository) CreateCaseUsedLog(mCase models.CaseUsedLog) (models.CaseUsedLog, error) {
	return mCase, r.db.DB.Create(&mCase).Error
}

func (r CaseUsedLogRepository) UpdateCaseUsedLog(caseId uuid.UUID, mCase models.CaseUsedLog) error {
	return r.db.DB.Model(&mCase).Where("case_id = ?", caseId).Updates(mCase).Error
}

func (r CaseUsedLogRepository) IncrementCaseUsedLog(caseId uuid.UUID) error {
	return r.db.DB.Model(&models.CaseUsedLog{}).Where("case_id = ?", caseId).Update("count", gorm.Expr("count + ?", 1)).Error
}
