package services

import (
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
	"mime/multipart"
)

type FileService struct {
	logger             *libs.Logger
	fileRepo           repositories.FileRepository
	fileTypeRepo       repositories.FileTypeRepository
	storageService     google_storage.GoogleStorage
	folderRepo         repositories.FolderRepository
	casePermissionRepo repositories.CasePermissionRepository
}

func NewFileService(
	logger *libs.Logger,
	fileRepo repositories.FileRepository,
	storageService google_storage.GoogleStorage,
	typeRepo repositories.FileTypeRepository,
	folderRepo repositories.FolderRepository,
	casePermissionRepo repositories.CasePermissionRepository,
) FileService {
	return FileService{
		logger:             logger,
		fileRepo:           fileRepo,
		fileTypeRepo:       typeRepo,
		storageService:     storageService,
		folderRepo:         folderRepo,
		casePermissionRepo: casePermissionRepo,
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

func (s *FileService) GetSignedUrl(amount int) ([]string, error) {
	return s.storageService.GetSignedUrl(amount)
}

func (s *FileService) UploadFile(
	file multipart.File,
	fileName string,
	fileType string,
	folderId uuid.UUID,
	userId uuid.UUID,
) (string, error) {

	err := s.checkPermission(userId, folderId)

	if err != nil {
		return "", err
	}

	_, err = s.storageService.UploadFile(file, fileName)

	if err != nil {
		return "", err
	}

	//fileT, err := s.fileTypeRepo.GetFileTypeByName(fileType)
	//
	//if err != nil {
	//	return "", err
	//}
	//
	//var fileUrl models.FileUrl
	//
	//fileUrl.Url, fileUrl.PreviewUrl = setFileURLs(fileType, url)
	//
	//modelFile := models.File{
	//	Name:     fileName,
	//	TypeId:   &fileT.ID,
	//	Urls:     []models.FileUrl{fileUrl},
	//	FolderId: &folderId,
	//}
	//
	//_, err = s.fileRepo.CreateFile(modelFile)
	//
	return "", nil
}

func setFileURLs(fileType string, url string) (string, string) {
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
			// TO DO : convert file to pdf and upload to google storage
			previewUrl := url + "&preview=true"
			return url, previewUrl
		}
	}

	return url, url
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

