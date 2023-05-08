package services

import (
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
	tagRepo            repositories.TagRepository
	actionRepo         repositories.ActionRepository
}

func NewFolderService(
	logger *libs.Logger,
	r repositories.FolderRepository,
	fileRepo repositories.FileRepository,
	casePermissionRepo repositories.CasePermissionRepository,
	storageService google_storage.GoogleStorage,
	casedUsedRepo repositories.CaseUsedLogRepository,
	actionLogRepo repositories.ActionLogRepository,
	tagRepo repositories.TagRepository,
	actionRepo repositories.ActionRepository,
) FolderService {
	return FolderService{
		logger:             logger,
		folderRepo:         r,
		fileRepo:           fileRepo,
		casePermissionRepo: casePermissionRepo,
		storageService:     storageService,
		casedUsedRepo:      casedUsedRepo,
		actionLogRepo:      actionLogRepo,
		tagRepo:            tagRepo,
		actionRepo:         actionRepo,
	}
}

func (s *FolderService) GetFolders() (folders []models.Folder, err error) {
	return s.folderRepo.GetFolders()
}

func (s *FolderService) GetFolder(id uuid.UUID, userId uuid.UUID) (dto *dtos.FolderDto, err error) {

	if err != nil {
		return nil, err
	}

	//err = s.casedUsedRepo.IncrementCaseUsedLog(*folder.CaseId, userId)

	folderModel, err := s.folderRepo.GetFolderContent(id)

	if err != nil {
		return nil, err
	}

	// FIXME : Fix this
	//_, err = s.casedUsedRepo.CreateCaseUsedLog(models.CaseUsedLog{CaseId: , UserId: userId})
	_, err = s.casedUsedRepo.FindOrCreate(*folderModel.CaseId, userId)

	if err != nil {
		return nil, err
	}

	err = s.casedUsedRepo.IncrementCaseUsedLog(*folderModel.CaseId, userId)

	if err != nil {
		return nil, err
	}

	urls, err := s.storageService.GetSignedFileUrls(folderModel.Files)

	files := folderModel.Files
	if err != nil {
		return nil, err
	}

	newFiles := make([]models.File, len(files))
	for i, file := range files {
		newFiles[i] = file
		newFiles[i].Url = &models.FileUrl{
			Url:        urls[i].Url,
			PreviewUrl: urls[i].PreviewUrl,
		}
	}

	s.logger.Info("Files: ", newFiles)

	folderModel.Files = newFiles

	folder := dtos.ToFolderDto(*folderModel)

	dto = &folder

	return dto, nil
}

func (s *FolderService) CreateFolder(dto dtos.CreateFolderDto, userId uuid.UUID) (*uuid.UUID, error) {

	parentFolderId, err := uuid.Parse(dto.ParentFolderId)

	if err != nil {
		return nil, err
	}

	parent, err := s.folderRepo.GetFolderContent(parentFolderId)

	if err != nil {
		return nil, err
	}

	tagFolder, err := s.tagRepo.GetTagByName("folder")

	if err != nil {
		return nil, err
	}

	folder := dto.ToModel(*parent.CaseId)

	folder.Tags = []models.Tag{*tagFolder}

	folder, err = s.folderRepo.CreateFolder(folder)

	if err != nil {
		return nil, err
	}

	err = s.addActionLog("create", folder.ID, userId)

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

	err = s.addActionLog("update", folder.ID, userId)

	return s.folderRepo.UpdateFolder(id, folder)
}

func (s *FolderService) MoveFolder(id uuid.UUID, dto dtos.MoveFolderDto, userId uuid.UUID) error {

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

func (s *FolderService) GetFolderLogs(userId, folderId uuid.UUID) ([]dtos.ActionLogDto, error) {

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
			urls, err := s.storageService.GetSignedFileUrls([]models.File{*log.File})
			if err != nil {
				s.logger.Info(err)
				return nil, err
			}
			url = urls[0]
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

func (s *FolderService) GetFolderRoot(userId, folderId uuid.UUID) ([]dtos.SimpleFolderDto, error) {

	err := s.checkPermission(userId, folderId)

	if err != nil {
		return nil, err
	}

	folders, err := s.folderRepo.GetFromRootToFolder(folderId)

	if err != nil {
		return nil, err
	}

	folderDtos := make([]dtos.SimpleFolderDto, len(folders))

	for i, folder := range folders {
		folderDtos[i] = dtos.SimpleFolderDto{
			Id:   folder.ID.String(),
			Name: folder.Name,
		}
	}

	return folderDtos, nil

}

func (s *FolderService) GetTagMenus(folderId uuid.UUID) ([]dtos.TagCountDto, error) {

	tags, err := s.tagRepo.CountFilesInTagsByFolderId(folderId)

	if err != nil {
		return nil, err
	}

	tagMenus := make([]dtos.TagCountDto, len(tags))

	for i, tag := range tags {
		tagMenus[i] = dtos.ToTagCountDto(tag)
	}

	return tagMenus, nil
}

func (s *FolderService) GetFileInTagId(folderId, tagId uuid.UUID) ([]dtos.FileDto, error) {

	files, err := s.fileRepo.GetFilesByFolderIdAndTagId(folderId, tagId)

	if err != nil {
		return nil, err
	}

	urls, err := s.storageService.GetSignedFileUrls(files)

	if err != nil {
		return nil, err
	}

	for i, _ := range files {
		files[i].Url = &models.FileUrl{
			Url:        urls[i].Url,
			PreviewUrl: urls[i].PreviewUrl,
		}
	}

	return dtos.ToFileDtos(files), nil
}

func (s *FolderService) addActionLog(actionName string, folderId, userId uuid.UUID) error {
	action, err := s.actionRepo.GetActionByName(actionName)

	if err != nil {
		return err
	}

	actionLog := models.ActionLog{
		FolderId: folderId,
		UserId:   userId,
		ActionId: action.ID,
	}

	if err != nil {
		return err
	}

	_, err = s.actionLogRepo.CreateActionLog(actionLog)

	return err
}
