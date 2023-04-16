package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type FileUrlRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFileUrlRepository(logger *libs.Logger, db libs.Database) FileUrlRepository {
	return FileUrlRepository{logger: logger, db: db}
}

func (r *FileUrlRepository) GetFileUrls() (fileUrls []models.FileUrl, err error) {
	return fileUrls, r.db.DB.Find(&fileUrls).Error
}

func (r *FileUrlRepository) GetFileUrl(id uuid.UUID) (fileUrl *models.FileUrl, err error) {
	return fileUrl, r.db.DB.First(&fileUrl, id).Error

}

func (r *FileUrlRepository) CreateFileUrl(fileUrl models.FileUrl) (models.FileUrl, error) {
	return fileUrl, r.db.DB.Create(&fileUrl).Error
}

func (r *FileUrlRepository) UpdateFileUrl(id uuid.UUID, fileUrl models.FileUrl) error {
	return r.db.DB.Model(&models.FileUrl{}).Where("id = ?", id).Updates(fileUrl).Error
}

func (r *FileUrlRepository) DeleteFileUrl(id uuid.UUID) error {
	return r.db.DB.Delete(&models.FileUrl{}, id).Error
}

func (r *FileUrlRepository) GetFileUrlByFileId(fileId uuid.UUID) (fileUrl *models.FileUrl, err error) {
	return fileUrl, r.db.DB.Where("file_id = ?", fileId).First(&fileUrl).Error
}
