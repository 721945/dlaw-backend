package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type CaseRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewCaseRepository(logger *libs.Logger, db libs.Database) CaseRepository {
	return CaseRepository{logger: logger, db: db}
}

func (r *CaseRepository) GetCases() (cases []models.Case, err error) {
	return cases, r.db.DB.Find(&cases).Preload("Folders").Error
}

func (r *CaseRepository) GetCase(id uuid.UUID) (mCase *models.Case, err error) {
	return mCase, r.db.DB.Preload("Folders").First(&mCase, id).Error
}

func (r *CaseRepository) CreateCase(mCase models.Case) (models.Case, error) {
	return mCase, r.db.DB.Create(&mCase).Error
}

func (r *CaseRepository) UpdateCase(id uuid.UUID, mCase models.Case) error {
	return r.db.DB.Model(&models.Case{}).Where("id = ?", id).Updates(mCase).Error
}

func (r *CaseRepository) DeleteCase(id uuid.UUID) error {
	return r.db.DB.Delete(&models.Case{}, id).Error
}
