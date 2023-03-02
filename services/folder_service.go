package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type FolderService struct {
	logger     *libs.Logger
	folderRepo repositories.FolderRepository
}

func NewFolderService(logger *libs.Logger, r repositories.FolderRepository) FolderService {
	return FolderService{logger: logger, folderRepo: r}
}

func (s *FolderService) GetFolders() (folders []models.Folder, err error) {
	return s.folderRepo.GetFolders()
}

func (s *FolderService) GetFolder(id uint) (folder *models.Folder, err error) {
	return s.folderRepo.GetFolder(id)
}

func (s *FolderService) CreateFolder(folder models.Folder) (models.Folder, error) {
	return s.folderRepo.CreateFolder(folder)
}

func (s *FolderService) UpdateFolder(id uint, folder models.Folder) error {
	return s.folderRepo.UpdateFolder(id, folder)
}

func (s *FolderService) DeleteFolder(id uint) error {
	return s.folderRepo.DeleteFolder(id)
}
