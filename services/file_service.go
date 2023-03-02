package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type FileService struct {
	logger   *libs.Logger
	fileRepo repositories.FileRepository
}

func NewFileService(logger *libs.Logger, r repositories.FileRepository) FileService {
	return FileService{logger: logger, fileRepo: r}
}

func (s *FileService) GetFiles() (files []models.File, err error) {
	return s.fileRepo.GetFiles()
}

func (s *FileService) GetFile(id uint) (file *models.File, err error) {
	return s.fileRepo.GetFile(id)
}

func (s *FileService) CreateFile(file models.File) (models.File, error) {
	return s.fileRepo.CreateFile(file)
}

func (s *FileService) UpdateFile(id uint, file models.File) error {
	return s.fileRepo.UpdateFile(id, file)
}

func (s *FileService) DeleteFile(id uint) error {
	return s.fileRepo.DeleteFile(id)
}
