package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type FileTypeRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFileTypeRepository(logger *libs.Logger, db libs.Database) FileTypeRepository {
	return FileTypeRepository{logger: logger, db: db}
}

func (r *FileTypeRepository) GetFileTypes() (fileTypes []models.FileType, err error) {
	return fileTypes, r.db.DB.Find(&fileTypes).Error
}

func (r *FileTypeRepository) GetFileType(id uuid.UUID) (fileType *models.FileType, err error) {
	return fileType, r.db.DB.First(&fileType, id).Error
}

func (r *FileTypeRepository) CreateFileType(fileType models.FileType) (models.FileType, error) {
	return fileType, r.db.DB.Create(&fileType).Error
}

func (r *FileTypeRepository) UpdateFileType(id uuid.UUID, fileType models.FileType) error {
	return r.db.DB.Model(&models.FileType{}).Where("id = ?", id).Updates(fileType).Error
}

func (r *FileTypeRepository) DeleteFileType(id uuid.UUID) error {
	return r.db.DB.Delete(&models.FileType{}, id).Error
}

func (r *FileTypeRepository) GetFileTypeByName(name string) (fileType *models.FileType, err error) {
	return fileType, r.db.DB.Where("name = ?", name).First(&fileType).Error
}
