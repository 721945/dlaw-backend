package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type FolderService struct {
	logger     *libs.Logger
	folderRepo repositories.FolderRepository
	fileRepo   repositories.FileRepository
}

func NewFolderService(logger *libs.Logger, r repositories.FolderRepository, fileRepo repositories.FileRepository) FolderService {
	return FolderService{logger: logger, folderRepo: r, fileRepo: fileRepo}
}

func (s *FolderService) GetFolders() (folders []models.Folder, err error) {
	return s.folderRepo.GetFolders()
}

func (s *FolderService) GetFolder(id uuid.UUID, userId uuid.UUID) (folder *models.Folder, err error) {
	folder, err = s.folderRepo.GetFolder(id)

	if err != nil {
		return nil, err
	}

	subFolders, err := s.folderRepo.GetSubFolders(id)

	s.logger.Info(subFolders)

	if err != nil {
		return nil, err
	}

	folder.SubFolders = subFolders

	files, err := s.fileRepo.GetFilesByFolderId(id)

	if err != nil {
		s.logger.Info(err)
		return nil, err
	}

	folder.Files = files

	return folder, nil
}

func (s *FolderService) CreateFolder(folder models.Folder) (*uuid.UUID, error) {
	folder, err := s.folderRepo.CreateFolder(folder)

	if err != nil {
		return nil, err
	}

	return &folder.ID, nil
}

func (s *FolderService) UpdateFolder(id uuid.UUID, folder models.Folder) error {
	return s.folderRepo.UpdateFolder(id, folder)
}

func (s *FolderService) DeleteFolder(id uuid.UUID) error {
	return s.folderRepo.DeleteFolder(id)
}
