package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type FileRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFileRepository(logger *libs.Logger, db libs.Database) FileRepository {
	return FileRepository{logger: logger, db: db}
}

func (r *FileRepository) GetFiles() (files []models.File, err error) {
	return files, r.db.DB.Find(&files).Error
}

func (r *FileRepository) GetFile(id uuid.UUID) (file *models.File, err error) {
	return file, r.db.DB.First(&file, id).Error
}

func (r *FileRepository) GetFilesByFolderId(folderId uuid.UUID) (files []models.File, err error) {
	return files, r.db.DB.Preload(clause.Associations).Where("folder_id = ?", folderId).Find(&files).Error
}

func (r *FileRepository) CreateFile(file models.File) (models.File, error) {
	return file, r.db.DB.Create(&file).Error
}

func (r *FileRepository) UpdateFile(id uuid.UUID, file models.File) error {
	return r.db.DB.Model(&models.File{}).Where("id = ?", id).Updates(file).Error
}

func (r *FileRepository) DeleteFile(id uuid.UUID) error {
	return r.db.DB.Delete(&models.File{}, id).Error
}
