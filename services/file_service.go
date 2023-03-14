package services

import (
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"mime/multipart"
)

type FileService struct {
	logger         *libs.Logger
	fileRepo       repositories.FileRepository
	storageService google_storage.GoogleStorage
}

func NewFileService(logger *libs.Logger, fileRepo repositories.FileRepository, storageService google_storage.GoogleStorage) FileService {
	return FileService{logger: logger, fileRepo: fileRepo, storageService: storageService}
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

func (s *FileService) GetSignedUrl(amount int) ([]string, error) {
	return s.storageService.GetSignedUrl(amount)
}

func (s *FileService) UploadFile(file multipart.File, fileName string) (string, error) {
	return s.storageService.UploadFile(file, fileName)
}
