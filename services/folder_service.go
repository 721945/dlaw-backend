package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type FolderService struct {
	logger             *libs.Logger
	folderRepo         repositories.FolderRepository
	fileRepo           repositories.FileRepository
	casePermissionRepo repositories.CasePermissionRepository
	storageService     google_storage.GoogleStorage
}

func NewFolderService(
	logger *libs.Logger,
	r repositories.FolderRepository,
	fileRepo repositories.FileRepository,
	casePermissionRepo repositories.CasePermissionRepository,
	storageService google_storage.GoogleStorage,
) FolderService {
	return FolderService{
		logger:             logger,
		folderRepo:         r,
		fileRepo:           fileRepo,
		casePermissionRepo: casePermissionRepo,
		storageService:     storageService,
	}
}

func (s *FolderService) GetFolders() (folders []models.Folder, err error) {
	return s.folderRepo.GetFolders()
}

func (s *FolderService) GetFolder(id uuid.UUID, userId uuid.UUID) (folder *models.Folder, err error) {
	folder, err = s.folderRepo.GetFolderContent(id)

	if err != nil {
		return nil, err
	}

	//subFolders, err := s.folderRepo.GetSubFolders(id)

	//if err != nil {
	//	return nil, err
	//}
	//
	//folder.SubFolders = subFolders

	//files, err := s.fileRepo.GetFilesByFolderId(id)
	files := folder.Files

	// Check if cloudname and previewcloudname is the same or not to avoid duplicate
	// If it the same just use the cloudname
	//for i, file := range files {
	//	if file.CloudName == file.PreviewCloudName {
	//		files[i].PreviewCloudName = ""
	//	}
	//}

	cloudNames := make([]string, len(files))
	previewCloudNames := make([]string, len(files))
	downloadNames := make([]string, len(files))

	for i, file := range files {
		cloudNames[i] = file.CloudName
		downloadNames[i] = file.Name
		previewCloudNames[i] = file.PreviewCloudName
	}

	s.logger.Info("cloudNames: ", cloudNames)
	s.logger.Info("previewCloudNames: ", previewCloudNames)
	s.logger.Info("downloadNames: ", downloadNames)
	//
	urls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)

	if err != nil {
		s.logger.Info(err)
		return nil, err
	}

	s.logger.Info("urls: ", urls)
	previewUrls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)
	if err != nil {
		s.logger.Info(err)
		return nil, err
	}
	s.logger.Info("previewUrls: ", previewUrls)

	newFiles := make([]models.File, len(files))
	for i, file := range files {
		newFiles[i] = file
		newFiles[i].Url = &models.FileUrl{
			Url:        urls[i],
			PreviewUrl: previewUrls[i],
		}
	}

	s.logger.Info("Files: ", newFiles)

	//if err != nil {
	//	s.logger.Info(err)
	//	return nil, err
	//}
	//

	folder.Files = newFiles

	return folder, nil
}

func (s *FolderService) CreateFolder(dto dtos.CreateFolderDto) (*uuid.UUID, error) {

	parentFolderId, err := uuid.Parse(dto.ParentFolderId)

	if err != nil {
		return nil, err
	}

	parent, err := s.folderRepo.GetFolderContent(parentFolderId)

	if err != nil {
		return nil, err
	}

	folder := dto.ToFolder(*parent.CaseId)

	folder, err = s.folderRepo.CreateFolder(folder)

	if err != nil {
		return nil, err
	}

	return &folder.ID, nil
}

func (s *FolderService) UpdateFolder(id uuid.UUID, dto dtos.UpdateFolderDto, userId uuid.UUID) error {

	err := s.checkPermission(userId, id)

	if err != nil {
		return err
	}

	folder := dto.ToFolder()

	return s.folderRepo.UpdateFolder(id, folder)
}

func (s *FolderService) DeleteFolder(id uuid.UUID, userId uuid.UUID) error {

	err := s.checkPermission(userId, id)

	if err != nil {
		return err
	}

	return s.folderRepo.DeleteFolder(id)
}

func (s *FolderService) ArchiveFolder(id uuid.UUID, userId uuid.UUID) error {

	err := s.checkPermission(userId, id)

	if err != nil {
		return err
	}

	return s.folderRepo.ArchiveFolder(id)
}

func (s *FolderService) UnArchiveFolder(id uuid.UUID, userId uuid.UUID) error {

	err := s.checkPermission(userId, id)

	if err != nil {
		return err
	}

	return s.folderRepo.UnArchiveFolder(id)
}

func (s *FolderService) checkPermission(userId uuid.UUID, folderId uuid.UUID) error {
	folder, err := s.folderRepo.GetFolder(folderId)

	if err != nil {
		return err
	}

	performerRole, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(userId, *folder.CaseId)

	if err != nil {
		return libs.ErrUnauthorized
	}

	permission := (*performerRole).Permission.Name

	if permission == "owner" || permission == "admin" {
		return nil
	}
	return libs.ErrUnauthorized
}
