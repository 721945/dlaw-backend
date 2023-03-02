package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
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

func (r *FileTypeRepository) GetFileType(id uint) (fileType *models.FileType, err error) {
	return fileType, r.db.DB.First(&fileType, id).Error
}

func (r *FileTypeRepository) CreateFileType(fileType models.FileType) (models.FileType, error) {
	return fileType, r.db.DB.Create(&fileType).Error
}

func (r *FileTypeRepository) UpdateFileType(id uint, fileType models.FileType) error {
	return r.db.DB.Model(&models.FileType{}).Where("id = ?", id).Updates(fileType).Error
}

func (r *FileTypeRepository) DeleteFileType(id uint) error {
	return r.db.DB.Delete(&models.FileType{}, id).Error
}
