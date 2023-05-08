package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type FolderRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFolderRepository(logger *libs.Logger, db libs.Database) FolderRepository {
	return FolderRepository{logger: logger, db: db}
}

func (r *FolderRepository) GetFolders() (folders []models.Folder, err error) {
	return folders, r.db.DB.Find(&folders).Error
}

func (r *FolderRepository) GetFolderContent(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.Preload("SubFolders").Preload("SubFolders.Tags").Preload("Files").Preload("Files.Tags").Preload("Tags").First(&folder, id).Error
}
func (r *FolderRepository) GetFolderContentWithOutFiles(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.Preload("SubFolders").Preload("SubFolders.Tags").Preload("Tags").First(&folder, id).Error
}

func (r *FolderRepository) GetFolder(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.First(&folder, id).Error
}

func (r *FolderRepository) GetRootFolder(caseId uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.Where("case_id = ? AND parent_folder_id IS NULL", caseId).First(&folder).Error
}

func (r *FolderRepository) GetFolderWithPermission(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.First(&folder, id).Error
}

func (r *FolderRepository) GetSubFolders(id uuid.UUID) (folders []models.Folder, err error) {
	return folders, r.db.DB.Where("parent_folder_id = ?", id).Find(&folders).Error
}

func (r *FolderRepository) CreateFolder(folder models.Folder) (models.Folder, error) {
	return folder, r.db.DB.Create(&folder).Error
}

func (r *FolderRepository) UpdateFolder(id uuid.UUID, folder models.Folder) error {
	return r.db.DB.Model(&models.Folder{}).Where("id = ?", id).Updates(folder).Error
}

func (r *FolderRepository) DeleteFolder(id uuid.UUID) error {
	return r.db.DB.Delete(&models.Folder{}, id).Error
}

func (r *FolderRepository) ArchiveFolder(id uuid.UUID) error {
	return r.db.DB.Model(&models.Folder{}).Where("id = ?", id).Update("is_archived", true).Error
}

func (r *FolderRepository) UnArchiveFolder(id uuid.UUID) error {
	return r.db.DB.Model(&models.Folder{}).Where("id = ?", id).Update("is_archived", false).Error
}

//func (r *FolderRepository) UpdateTags(id uuid.UUID, tags []models.Tag) error {
//	//err := r.db.DB.Model(&models.Folder{}).Where("id = ?", id).Association("Tags").Clear()
//	//if err != nil {
//	//	r.logger.Info(err)
//	//	return err
//	//}
//return r.db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Model(&models.Folder{}).Where("id = ?", id).Updates(models.Folder{Tags: tags}).Error
//}

func (r *FolderRepository) UpdateTags(id uuid.UUID, tags []models.Tag) error {
	model := models.Folder{
		Base: models.Base{ID: id},
	}
	err := r.db.DB.Model(&model).Association("Tags").Clear()
	if err != nil {
		r.logger.Info(err)
		return err
	}
	return r.db.DB.Model(&model).Association("Tags").Append(&tags)
}

func (r *FolderRepository) GetParentFolder(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.Where("id = ?", id).First(&folder).Error
}

func (r *FolderRepository) GetFromRootToFolder(id uuid.UUID) (folders []models.Folder, err error) {

	folder, err := r.GetParentFolder(id)

	//
	if folder.ParentFolderId != nil {

		parentFolder, err := r.GetFromRootToFolder(*folder.ParentFolderId)
		if err != nil {
			return folders, err
		}
		folders = append(folders, parentFolder...)
	}

	folders = append(folders, *folder)
	//
	return folders, err
}

// GetFoldersByCaseId returns all folders for a case
func (r *FolderRepository) GetFoldersByCaseId(caseId uuid.UUID) (folders []models.Folder, err error) {
	return folders, r.db.DB.Where("case_id = ?", caseId).Find(&folders).Error
}

//func (r *FolderRepository) GetCasePermissionByFolderId(folderId uuid.UUID) (folders []models.Folder, err error) {
//
//}
