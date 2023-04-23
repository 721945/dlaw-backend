package services

import (
	"errors"
	"fmt"
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
	casedUsedRepo      repositories.CaseUsedLogRepository
	actionLogRepo      repositories.ActionLogRepository
}

func NewFolderService(
	logger *libs.Logger,
	r repositories.FolderRepository,
	fileRepo repositories.FileRepository,
	casePermissionRepo repositories.CasePermissionRepository,
	storageService google_storage.GoogleStorage,
	casedUsedRepo repositories.CaseUsedLogRepository,
	actionLogRepo repositories.ActionLogRepository,
) FolderService {
	return FolderService{
		logger:             logger,
		folderRepo:         r,
		fileRepo:           fileRepo,
		casePermissionRepo: casePermissionRepo,
		storageService:     storageService,
		casedUsedRepo:      casedUsedRepo,
		actionLogRepo:      actionLogRepo,
	}
}

func (s *FolderService) GetFolders() (folders []models.Folder, err error) {
	return s.folderRepo.GetFolders()
}

func (s *FolderService) GetFolder(id uuid.UUID, userId uuid.UUID) (folder *models.Folder, err error) {
	folder, err = s.folderRepo.GetFolderContent(id)

	err = s.casedUsedRepo.IncrementCaseUsedLog(*folder.CaseId)

	if err != nil {
		return nil, err
	}

	files := folder.Files

	cloudNames := make([]string, len(files))
	previewCloudNames := make([]string, len(files))
	downloadNames := make([]string, len(files))

	for i, file := range files {
		cloudNames[i] = file.CloudName
		downloadNames[i] = file.Name
		previewCloudNames[i] = file.PreviewCloudName
	}

	urls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)

	if err != nil {
		s.logger.Info(err)
		return nil, err
	}

	previewUrls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)
	if err != nil {
		s.logger.Info(err)
		return nil, err
	}

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

	folder := dto.ToModel(*parent.CaseId)

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

	checkFolder, err := s.folderRepo.GetFolder(id)

	if err != nil {
		return err
	}

	if checkFolder.ParentFolderId == nil {
		return fmt.Errorf("Can not update root folder")
	}

	folder := dto.ToModel()

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

// GetActionLogs Get action logs
func (s *FolderService) GetActionLogs(folderId uuid.UUID, userId uuid.UUID) ([]dtos.ActionLogDto, error) {

	err := s.checkPermission(userId, folderId)

	if err != nil {
		return nil, err
	}

	logs, err := s.actionLogRepo.GetActionLogsByFolderId(folderId)

	if err != nil {
		return nil, err
	}

	actionLogs := make([]dtos.ActionLogDto, len(logs))

	for i, log := range logs {
		var url models.FileUrl
		if log.File != nil {
			url, err = s.getSignedFileUrl(*log.File)
			if err != nil {
				s.logger.Info(err)
				return nil, err
			}
		}
		urlDto := dtos.ToFileActionLogDto(log.File, url.PreviewUrl)
		actionLogs[i] = dtos.ToActionLogDto(log, urlDto)
	}

	return actionLogs, nil
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

// get signed url for files
// TODO : Refactor this to option use download or not
func (s *FolderService) getSignedFileUrls(files []models.File) (fileUrls []models.FileUrl, err error) {

	cloudNames := make([]string, len(files))
	previewCloudNames := make([]string, len(files))
	downloadNames := make([]string, len(files))

	for i, file := range files {
		cloudNames[i] = file.CloudName
		downloadNames[i] = file.Name
		previewCloudNames[i] = file.PreviewCloudName
	}

	urlsCh := make(chan []string)
	previewUrlsCh := make(chan []string)

	// Run the two calls to GetSignedUrls in parallel using goroutines
	go func() {
		urls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)
		if err != nil {
			s.logger.Info(err)
			urlsCh <- nil
		} else {
			urlsCh <- urls
		}
	}()

	go func() {
		previewUrls, err := s.storageService.GetSignedUrls(previewCloudNames, []string{}, downloadNames)
		if err != nil {
			s.logger.Info(err)
			previewUrlsCh <- nil
		} else {
			previewUrlsCh <- previewUrls
		}
	}()

	// Wait for both goroutines to complete and merge the results
	urls := <-urlsCh
	previewUrls := <-previewUrlsCh

	if urls == nil || previewUrls == nil {
		return nil, errors.New("error getting signed urls")
	}

	fileUrls = make([]models.FileUrl, len(files))

	for i, _ := range files {
		fileUrls[i] = models.FileUrl{
			Url:        urls[i],
			PreviewUrl: previewUrls[i],
		}
	}

	return fileUrls, nil
}

// getSignedFileUrl
func (s *FolderService) getSignedFileUrl(file models.File) (fileUrl models.FileUrl, err error) {

	cloudNames := []string{file.CloudName}
	previewCloudNames := []string{file.PreviewCloudName}
	downloadNames := []string{file.Name}

	urls, err := s.storageService.GetSignedUrls(cloudNames, []string{}, downloadNames)
	if err != nil {
		s.logger.Info(err)
		return models.FileUrl{}, err
	}

	previewUrls, err := s.storageService.GetSignedUrls(previewCloudNames, []string{}, downloadNames)
	if err != nil {
		s.logger.Info(err)
		return models.FileUrl{}, err
	}

	fileUrl = models.FileUrl{
		Url:        urls[0],
		PreviewUrl: previewUrls[0],
	}

	return fileUrl, nil
}
