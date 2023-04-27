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

func (r CaseUsedLogRepository) IncrementCaseUsedLog(caseId uuid.UUID, userId uuid.UUID) error {
	return r.db.DB.Model(&models.CaseUsedLog{}).Where("case_id = ? AND user_id = ?", caseId, userId).Update("count", gorm.Expr("count + ?", 1)).Error
}

func (r CaseUsedLogRepository) FindOrCreate(caseId uuid.UUID, userId uuid.UUID) (models.CaseUsedLog, error) {
	var mCase models.CaseUsedLog
	mCase = models.CaseUsedLog{CaseId: caseId, UserId: userId}
	err := r.db.DB.Where("case_id = ? AND user_id = ?", caseId, userId).FirstOrCreate(&mCase).Error
	return mCase, err
}

func (r CaseUsedLogRepository) GetCaseUsedLogWithCaseByCaseIdsAndUserId(caseIds []uuid.UUID, userId uuid.UUID) (cases []models.CaseUsedLog, err error) {
	return cases, r.db.DB.Preload("Case", "is_archive = false").Preload("Case.Folders", "parent_folder_id IS NULL").Where("case_id IN (?) AND user_id = ?", caseIds, userId).Find(&cases).Error
}
