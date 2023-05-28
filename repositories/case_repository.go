package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type CaseRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewCaseRepository(logger *libs.Logger, db libs.Database) CaseRepository {
	return CaseRepository{logger: logger, db: db}
}

func (r *CaseRepository) GetCases() (cases []models.Case, err error) {
	return cases, r.db.DB.Preload("Folders").Find(&cases).Error
}

func (r *CaseRepository) GetCase(id uuid.UUID) (mCase *models.Case, err error) {
	return mCase, r.db.DB.Preload("Folders", "parent_folder_id IS NULL").First(&mCase, id).Error
}

func (r *CaseRepository) CreateCase(mCase models.Case) (models.Case, error) {
	return mCase, r.db.DB.Create(&mCase).Error
}

func (r *CaseRepository) UpdateCase(id uuid.UUID, mCase models.Case) error {
	return r.db.DB.Model(&models.Case{}).Where("id = ?", id).Updates(mCase).Error
}

func (r *CaseRepository) DeleteCase(id uuid.UUID) error {
	return r.db.DB.Select(clause.Associations).Delete(&models.Case{}, id).Error
}

func (r *CaseRepository) GetCasesByIds(caseIds []uuid.UUID, isArchive bool) (cases []models.Case, err error) {
	return cases, r.db.DB.Preload("Folders", "parent_folder_id IS NULL").Preload("CasePermissions").Preload("CasePermissions.User").Preload("CasePermissions.Permission").Where("id IN (?) AND is_archive = ?", caseIds, isArchive).Find(&cases).Error
}

func (r *CaseRepository) GetCasesByFolderId(folderId uuid.UUID) (cases []models.Case, err error) {
	return cases, r.db.DB.Preload("Folders", "parent_folder_id = ?", folderId).Find(&cases).Error
}

func (r *CaseRepository) ArchiveCase(id uuid.UUID) error {
	return r.db.DB.Model(&models.Case{}).Where("id = ?", id).Update("is_archive", true).Error
}

func (r *CaseRepository) UnArchiveCase(id uuid.UUID) error {
	return r.db.DB.Model(&models.Case{}).Where("id = ?", id).Update("is_archive", false).Error
}

func (r *CaseRepository) GetCasesSortedByFrequency(userId uuid.UUID) (cases []models.Case, err error) {
	//return cases, r.db.DB.Preload("CaseUsedLogs", "user_id = ?", userId).Find(&cases).Error
	//return cases, r.db.DB.Joins("JOIN case_used_logs ON case_used_logs.user_id = users.id").Preload("Folders", "parent_folder_id IS NULL").Order("count DESC").Limit(10).Find(&cases).Error
	return cases, r.db.DB.Preload("CaseUsedLogs", "user_id = ?", userId).Preload("Folders", "parent_folder_id IS NULL").Preload("Folders.Tags").Order("count DESC").Limit(10).Find(&cases).Error
}

func (r *CaseRepository) GetCasesWhichFileIsPublic() (cases []models.Case, err error) {
	subQuery := r.db.DB.Table("files").Select("1").Where("files.case_id = cases.id AND files.is_public = ? AND files.deleted_at IS NULL", true)

	return cases, r.db.DB.Preload("Files", "is_public = true").Preload("Files.FileType").Where("EXISTS (?)", subQuery).Find(&cases).Error
}

//func (r *CaseRepository) GetCasesWhichFileIsPublic() (cases []models.Case, err error) {
//	return cases, r.db.DB.Preload("Files", "is_public = true").Find(&cases).Error
//}
