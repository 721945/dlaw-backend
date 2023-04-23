package services

import (
	"fmt"
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
	}
}

func (s *FileService) GetFiles() (files []models.File, err error) {
	return s.fileRepo.GetFiles()
}

func (s *FileService) GetFile(id uuid.UUID) (file *models.File, err error) {
	return s.fileRepo.GetFile(id)
}

func (s *FileService) CreateFile(file models.File) (models.File, error) {
	return s.fileRepo.CreateFile(file)
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

	err := s.checkPermission(userId, folderId)
	mimeTypeSTtring := convertMimeTypeToString(fileType)
	if err != nil {
		return "", err
	}
	fileT, err := s.fileTypeRepo.GetFileTypeByName(mimeTypeSTtring)

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

	tag, err := s.tagRepo.GetTagByNames([]string{mimeTypeSTtring})

	// Do the ocr and then update the tag

	modelFile := models.File{
		Name:             fileName,
		TypeId:           &fileT.ID,
		FolderId:         &folderId,
		CloudName:        cloudName,
		PreviewCloudName: previewCloudName,
		Tags:             tag,
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

func (s *FileService) checkPermission(userId uuid.UUID, folderId uuid.UUID) error {
	// TODO : Improve to get permission by folder id later
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
