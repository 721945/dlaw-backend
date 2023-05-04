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
	"mime/multipart"
	"strings"
	"time"
)

type FileService struct {
	logger             *libs.Logger
	fileRepo           repositories.FileRepository
	fileTypeRepo       repositories.FileTypeRepository
	storageService     google_storage.GoogleStorage
	folderRepo         repositories.FolderRepository
	casePermissionRepo repositories.CasePermissionRepository
	actionRepo         repositories.ActionRepository
	actionLogRepo      repositories.ActionLogRepository
	tagRepo            repositories.TagRepository
	fileViewLogRepo    repositories.FileViewLogRepository
}

func NewFileService(
	logger *libs.Logger,
	fileRepo repositories.FileRepository,
	storageService google_storage.GoogleStorage,
	typeRepo repositories.FileTypeRepository,
	folderRepo repositories.FolderRepository,
	casePermissionRepo repositories.CasePermissionRepository,
	actionRepo repositories.ActionRepository,
	actionLogRepo repositories.ActionLogRepository,
	tagRepo repositories.TagRepository,
	fileViewLogRepo repositories.FileViewLogRepository,
) FileService {
	return FileService{
		logger:             logger,
		fileRepo:           fileRepo,
		fileTypeRepo:       typeRepo,
		storageService:     storageService,
		folderRepo:         folderRepo,
		casePermissionRepo: casePermissionRepo,
		actionRepo:         actionRepo,
		actionLogRepo:      actionLogRepo,
		tagRepo:            tagRepo,
		fileViewLogRepo:    fileViewLogRepo,
	}
}

func (s *FileService) GetFiles() (dto []dtos.FileDto, err error) {
	files, err := s.fileRepo.GetFiles()

	if err != nil {
		return nil, err
	}

	urls, err := s.getSignedFileUrls(files)

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
	dto = make([]dtos.FileDto, 0)

	for _, file := range newFiles {
		dto = append(dto, dtos.ToFileDto(file))
	}

	return dto, err
}

func (s *FileService) GetFile(id uuid.UUID, userId *uuid.UUID) (dto *dtos.FileDto, err error) {
	file, err := s.fileRepo.GetFile(id)

	if userId != nil {
		_, err := s.fileViewLogRepo.CreateFileViewLog(models.FileViewLog{
			FileId: id,
			UserId: *userId,
		})

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	url, err := s.getSignedFileUrls([]models.File{*file})

	if err != nil {
		return nil, err
	}

	file.Url = &models.FileUrl{
		Url:        url[0].Url,
		PreviewUrl: url[0].PreviewUrl,
	}

	fileDto := dtos.ToFileDto(*file)

	return &fileDto, err
}

func (s *FileService) CreateFile(file models.File) (string, error) {
	f, err := s.fileRepo.CreateFile(file)

	if err != nil {
		return "", err
	}

	return f.ID.String(), nil
}

func (s *FileService) UpdateFile(id uuid.UUID, file models.File) error {
	return s.fileRepo.UpdateFile(id, file)
}

func (s *FileService) DeleteFile(id uuid.UUID) error {
	return s.fileRepo.DeleteFile(id)
}

//func (s *FileService) GetSignedUrl(amount int) ([]string, error) {
//	return s.storageService.GetSignedUrls(amount)
//}

func (s *FileService) UploadFile(
	file multipart.File,
	fileName string,
	fileType string,
	folderId uuid.UUID,
	userId uuid.UUID,
) (string, error) {
	// TODO: NEED TO CHECK FOR REPLACE FILE

	caseId, err := s.checkPermissionAndGetCaseId(userId, folderId)
	mimeTypeToString := convertMimeTypeToString(fileType)
	if err != nil {
		return "", err
	}
	fileT, err := s.fileTypeRepo.GetFileTypeByName(mimeTypeToString)

	if err != nil {
		return "", err
	}

	needConvert := checkNeedConvert(fileType)

	version, versionPreview := "", ""

	cloudName := generateUniqueName()

	previewCloudName := cloudName

	version, err = s.storageService.UploadFile(file, cloudName)

	versionPreview = version

	if needConvert {
		//	DO CONVERT
		//	previewCloudName = generateUniqueName(fileName)
		//	versionPreview, err = s.storageService.UploadFile(file, previewCloudName)
	}

	if err != nil {
		return "", err
	}

	tag, err := s.tagRepo.GetTagByNames([]string{mimeTypeToString})

	// Do the ocr and then update the tag

	modelFile := models.File{
		Name:             fileName,
		TypeId:           &fileT.ID,
		FolderId:         &folderId,
		CloudName:        cloudName,
		PreviewCloudName: previewCloudName,
		Tags:             tag,
		CaseId:           caseId,
	}

	//
	fileRes, err := s.fileRepo.CreateFile(modelFile)

	if err != nil {
		return "", err
	}

	action, err := s.actionRepo.GetActionByName("upload")

	if err != nil {
		return "", err
	}

	actionLog := models.ActionLog{
		FolderId:             folderId,
		FileId:               &fileRes.ID,
		UserId:               userId,
		ActionId:             action.ID,
		FileVersionId:        &version,
		FilePreviewVersionId: &versionPreview,
	}

	err = s.updateFolderTagByFolderId(folderId)

	if err != nil {
		return "", err
	}

	_, err = s.actionLogRepo.CreateActionLog(actionLog)

	if err != nil {
		return "", err
	}

	return fileRes.ID.String(), nil
}

func generateUniqueName() string {
	id := uuid.New().String()
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%s-%d", id, timestamp)
}

func checkNeedConvert(fileType string) bool {
	microsoftTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/msword",
	}

	for _, t := range microsoftTypes {
		if t == fileType {
			return true
		}
	}

	return false
}

func (s *FileService) checkPermissionAndGetCaseId(userId uuid.UUID, folderId uuid.UUID) (*uuid.UUID, error) {
	// TODO : Improve to get permission by folder id later
	folder, err := s.folderRepo.GetFolder(folderId)

	if err != nil {
		return nil, err
	}

	performerRole, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(userId, *folder.CaseId)

	if err != nil {
		return nil, libs.ErrUnauthorized
	}

	permission := (*performerRole).Permission.Name

	if permission == "owner" || permission == "admin" {
		return folder.CaseId, nil
	}
	return nil, libs.ErrUnauthorized
}

func convertMimeTypeToString(mimeType string) string {
	//if start with image

	if strings.HasPrefix(mimeType, "image") {
		return "image"
	}

	if strings.HasPrefix(mimeType, "video") {
		return "video"
	}

	//audio
	if strings.HasPrefix(mimeType, "audio") {
		return "audio"
	}

	// compress
	if strings.HasPrefix(mimeType, "application/zip") {
		return "compress"
	}

	// text
	if strings.HasPrefix(mimeType, "text") {
		return "text"
	}

	if strings.HasPrefix(mimeType, "audio") {
		return "audio"
	}

	switch mimeType {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "word"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "excel"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "powerpoint"
	case "application/vnd.ms-excel":
		return "excel"
	case "application/vnd.ms-powerpoint":
		return "powerpoint"
	case "application/msword":
		return "word"
	case "application/pdf":
		return "pdf"
	case "image/jpeg":
		return "image"
	case "image/png":
		return "image"
	case "image/gif":
		return "image"
	default:
		return ""
	}
}

func (s *FileService) updateFolderTagByFolderId(folderId uuid.UUID) error {
	// Get all tags from files in folder
	files, err := s.fileRepo.GetFilesWithTagByFolderId(folderId)

	if err != nil {
		return err
	}

	tag, err := s.tagRepo.GetTagByNames([]string{"folder"})

	if err != nil {
		return err
	}
	// make it distinct
	tags := make(map[string]models.Tag)

	tags[tag[0].Name] = tag[0]

	for _, file := range files {
		for _, tag := range file.Tags {
			tags[tag.Name] = tag
		}
	}

	tagArray := make([]models.Tag, 0)

	for _, tag := range tags {
		tagArray = append(tagArray, tag)
	}

	err = s.folderRepo.UpdateTags(folderId, tagArray)

	return err

}

func (s *FileService) getSignedFileUrls(files []models.File) (fileUrls []models.FileUrl, err error) {

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

func (s *FileService) CountFilesInTags(userId uuid.UUID) ([]dtos.TagCountDto, error) {

	casePermissions, err := s.casePermissionRepo.GetCasePermissionsByUserId(userId)

	if err != nil {
		return nil, err
	}

	caseIds := make([]uuid.UUID, len(casePermissions))

	for i, casePermission := range casePermissions {
		caseIds[i] = casePermission.CaseId
	}

	files, err := s.fileRepo.GetFilesByCaseIds(caseIds)

	if err != nil {
		return nil, err
	}

	fileIds := make([]uuid.UUID, len(files))

	for i, file := range files {
		fileIds[i] = file.ID
	}

	tagCount, err := s.tagRepo.CountFilesInTags(fileIds)

	s.logger.Info(tagCount)

	if err != nil {
		return nil, err
	}

	tagDtos := make([]dtos.TagCountDto, len(tagCount))

	for i, tag := range tagCount {
		tagDtos[i] = dtos.ToTagCountDto(tag)
	}

	return tagDtos, nil
}

func (s *FileService) GetRecentFileOpened(userId uuid.UUID) ([]dtos.FileDto, error) {

	fileViews, err := s.fileViewLogRepo.GetFileViewLogsForUser(userId)

	if err != nil {
		return nil, err
	}

	files := make([]models.File, len(fileViews))

	for i, fileView := range fileViews {
		files[i] = *fileView.File
	}

	fileUrls, err := s.getSignedFileUrls(files)

	if err != nil {
		return nil, err
	}

	dto := make([]dtos.FileDto, len(fileViews))

	for i, fileView := range fileViews {
		fileView.File.Url = &fileUrls[i]
		dto[i] = dtos.ToFileDto(*fileView.File)
	}

	return dto, nil
}