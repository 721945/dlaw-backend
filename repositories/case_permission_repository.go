package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type CasePermissionRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewCasePermissionRepository(logger *libs.Logger, db libs.Database) CasePermissionRepository {
	return CasePermissionRepository{logger: logger, db: db}
}

func (r CasePermissionRepository) GetCasePermissions() (cases []models.CasePermission, err error) {
	return cases, r.db.DB.Find(&cases).Error
}

func (r CasePermissionRepository) GetCasePermission(id uuid.UUID) (mCase *models.CasePermission, err error) {
	return mCase, r.db.DB.First(&mCase, id).Error
}

func (r CasePermissionRepository) CreateCasePermission(mCase models.CasePermission) (models.CasePermission, error) {
	return mCase, r.db.DB.Create(&mCase).Error
}

func (r CasePermissionRepository) UpdateCasePermission(id uuid.UUID, mCase models.CasePermission) error {
	return r.db.DB.Model(&mCase).Where("id = ?", id).Updates(mCase).Error
}

func (r CasePermissionRepository) DeleteCasePermission(id uuid.UUID) error {
	return r.db.DB.Delete(&models.CasePermission{}, id).Error
}

func (r CasePermissionRepository) GetCasePermissionsByCaseId(id uuid.UUID) (cases []models.CasePermission, err error) {
	return cases, r.db.DB.Preload("User").Where("case_id = ?", id).Find(&cases).Error
}

func (r CasePermissionRepository) GetCasePermissionsByUserId(id uuid.UUID) (cases []models.CasePermission, err error) {
	return cases, r.db.DB.Preload("User").Where("user_id = ?", id).Find(&cases).Error
}

func (r CasePermissionRepository) GetCasePermissionsByUserIdAndCaseId(userId uuid.UUID, caseId uuid.UUID) (mCase *models.CasePermission, err error) {
	return mCase, r.db.DB.Preload("Permission").Where("user_id = ? AND case_id = ?", userId, caseId).First(&mCase).Error
}

//func (r CasePermissionRepository) GetCasePermissionByFolderId(folderId, userId uuid.UUID) (mCase *models.CasePermission, err error) {
//	return mCase, r.db.DB.Joins("Case").Joins("Folder").Where("folder.id = ? AND case_permission.user_id = ?", folderId, userId).Preload("Permission").First(&mCase).Error
//}
