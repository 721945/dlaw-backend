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

func (r *FolderRepository) GetFolder(id uuid.UUID) (folder *models.Folder, err error) {
	return folder, r.db.DB.Preload("Tags").First(&folder, id).Error
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
