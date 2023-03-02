package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
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

func (r *FileRepository) GetFile(id uint) (file *models.File, err error) {
	return file, r.db.DB.First(&file, id).Error
}

func (r *FileRepository) CreateFile(file models.File) (models.File, error) {
	return file, r.db.DB.Create(&file).Error
}

func (r *FileRepository) UpdateFile(id uint, file models.File) error {
	return r.db.DB.Model(&models.File{}).Where("id = ?", id).Updates(file).Error
}

func (r *FileRepository) DeleteFile(id uint) error {
	return r.db.DB.Delete(&models.File{}, id).Error
}
