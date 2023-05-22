package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type FileService struct {
	logger             *libs.Logger
	fileRepo           repositories.FileRepository
	fileTypeRepo       repositories.FileTypeRepository
	folderRepo         repositories.FolderRepository
	casePermissionRepo repositories.CasePermissionRepository
	actionRepo         repositories.ActionRepository
	actionLogRepo      repositories.ActionLogRepository
	tagRepo            repositories.TagRepository
	fileViewLogRepo    repositories.FileViewLogRepository
	storageService     google_storage.GoogleStorage
	ocrService         google_vision.GoogleVision
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
	ocrService google_vision.GoogleVision,
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
		ocrService:         ocrService,
	}
}

func (s *FileService) GetFiles() (dto []dtos.FileDto, err error) {
	files, err := s.fileRepo.GetFiles()

	if err != nil {
		return nil, err
	}

	urls, err := s.storageService.GetSignedFileUrls(files)

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
	file, err := s.fileRepo.GetFileContent(id)

	if err != nil {
		return nil, err
	}

	if !file.IsPublic && !file.IsShared {
		if userId == nil {
			return nil, libs.ErrNotFound
		}
	}

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

	url, err := s.storageService.GetSignedFileUrls([]models.File{*file})

	if err != nil {
		return nil, err
	}

	file.Url = &models.FileUrl{
		Url:        url[0].Url,
		PreviewUrl: url[0].PreviewUrl,
	}

	size, err := s.storageService.GetFileSize(file.CloudName, "")

	if err != nil {
		return nil, err
	}

	fileDto := dtos.ToFileWithSizeDto(*file, size)

	return &fileDto, err
}

func (s *FileService) CreateFile(file models.File) (string, error) {
	f, err := s.fileRepo.CreateFile(file)

	if err != nil {
		return "", err
	}

	return f.ID.String(), nil
}

func (s *FileService) UpdateFile(id uuid.UUID, dto dtos.UpdateFileDto, userId uuid.UUID) error {
	file, err := s.fileRepo.GetFile(id)

	if err != nil {
		return err
	}

	err = s.addLogs("rename", userId, *file.FolderId, id, "", "", &file.Name, &dto.Name)

	file.Name = dto.Name

	return s.fileRepo.UpdateFile(id, *file)
}

func (s *FileService) MoveFile(id uuid.UUID, dto dtos.MoveFileDto, userId uuid.UUID) error {
	file, err := s.fileRepo.GetFile(id)

	if err != nil {
		return err
	}

	folder, err := s.folderRepo.GetFolder(*file.FolderId)

	if err != nil {
		return err
	}

	toFolderId, err := uuid.Parse(dto.TargetFolderId)

	if err != nil {
		return err
	}

	toFolder, err := s.folderRepo.GetFolder(toFolderId)

	err = s.addLogs("move", userId, *file.FolderId, id, "", "", &folder.Name, &toFolder.Name)

	model := dto.ToModel()

	if model == nil {
		return errors.New("invalid dto")
	}

	return s.fileRepo.UpdateFile(id, *model)
}

func (s *FileService) ShareFile(id string, userId uuid.UUID) (string, error) {
	fileId, err := uuid.Parse(id)

	if err != nil {
		return "", err
	}

	file, err := s.fileRepo.GetFile(fileId)

	if err != nil {
		return "", err
	}

	links, err := s.storageService.GiveAccessPublic(file.CloudName, file.Name)

	isPublic := file.IsPublic
	isShared := file.IsShared

	fromText := ""

	if isPublic {
		fromText = "public"
	} else if isShared {
		fromText = "have-link"
	} else {
		fromText = "private"
	}

	toText := "private"

	err = s.addLogs("share", userId, *file.FolderId, fileId, "", "", &fromText, &toText)

	if err != nil {
		return "", err
	}

	err = s.fileRepo.UpdateFilePublic(fileId, models.File{
		IsShared: true,
		IsPublic: false,
	})

	if err != nil {
		return "", err
	}

	return links, nil
}

func (s *FileService) RemoveShareFile(id string, userId uuid.UUID) error {
	fileId, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	file, err := s.fileRepo.GetFile(fileId)

	if err != nil {
		return err
	}

	err = s.storageService.GiveAccessPrivate(file.CloudName)

	isPublic := file.IsPublic
	isShared := file.IsShared

	fromText := ""

	if isPublic {
		fromText = "public"
	} else if isShared {
		fromText = "have-link"
	} else {
		fromText = "private"
	}

	toText := "private"

	err = s.addLogs("share", userId, *file.FolderId, fileId, "", "", &fromText, &toText)

	if err != nil {
		return err
	}

	err = s.fileRepo.UpdateFilePublic(fileId, models.File{
		IsShared: false,
		IsPublic: false,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *FileService) PublicFile(id string, userId uuid.UUID) (string, error) {
	fileId, err := uuid.Parse(id)

	if err != nil {
		return "", err
	}

	file, err := s.fileRepo.GetFile(fileId)

	if err != nil {
		return "", err
	}

	links, err := s.storageService.GiveAccessPublic(file.CloudName, file.Name)

	if err != nil {
		return "", err
	}

	err = s.fileRepo.UpdateFilePublic(fileId, models.File{
		IsShared: true,
		IsPublic: true,
	})

	if err != nil {
		return "", err
	}

	isPublic := file.IsPublic
	isShared := file.IsShared

	fromText := ""

	if isPublic {
		fromText = "public"
	} else if isShared {
		fromText = "have-link"
	} else {
		fromText = "private"
	}

	toText := "public"

	err = s.addLogs("share", userId, *file.FolderId, fileId, "", "", &fromText, &toText)

	return links, nil
}

func (s *FileService) SearchFiles(word, caseId, folderId, tag, fileType, page, limit string, userID uuid.UUID) ([]dtos.FileDto, dtos.PaginationResponse, error) {
	var pagination dtos.PaginationResponse

	word = utils.ConvertThaiNumToEngNum(word)
	
	if word == "" && caseId == "" && folderId == "" && tag == "" && fileType == "" {
		return nil, pagination, errors.New("invalid search params")
	}

	var filters []string

	if caseId != "" {
		_, err := uuid.Parse(caseId)
		if err != nil {
			return nil, pagination, errors.New("invalid case id")
		}
		filters = append(filters, "case_id = \""+caseId+"\"")
	}

	if folderId != "" {
		_, err := uuid.Parse(folderId)
		if err != nil {
			return nil, pagination, errors.New("invalid folder id")
		}

		filters = append(filters, "folder_ids = \""+folderId+"\"")
	}

	if tag != "" {
		filters = append(filters, "tags = \""+tag+"\"")
	}

	if fileType != "" {
		filters = append(filters, "type = \""+fileType+"\"")
	}

	pageNum, err := strconv.ParseInt(page, 10, 64)

	if err != nil {
		return nil, pagination, errors.New("invalid page number")
	}

	limitNum, err := strconv.ParseInt(limit, 10, 64)

	if err != nil {
		return nil, pagination, errors.New("invalid limit number")
	}
	filter := strings.Join(filters, " AND ")
	var searchResult *meilisearch.SearchResponse

	searchResult, err = s.fileRepo.SearchFiles(word, filter, pageNum, limitNum)

	if err != nil {
		return nil, pagination, err
	}

	var ids []uuid.UUID
	for _, hit := range searchResult.Hits {
		var result models.MeiliFileResponse
		s.logger.Info("Hit: %v\n", hit)
		documentsBytes, _ := json.Marshal(hit)
		err := json.Unmarshal(documentsBytes, &result)
		if err != nil {
			s.logger.Error("Error: %v\n", err)
		}
		id, _ := uuid.Parse(result.ID)
		ids = append(ids, id)
	}

	files, err := s.fileRepo.GetFilesByIds(ids)

	urls, err := s.storageService.GetSignedFileUrls(files)

	if err != nil {
		return nil, pagination, err
	}

	newFiles := make([]models.File, len(files))

	for i, file := range files {
		newFiles[i] = file
		newFiles[i].Url = &models.FileUrl{
			Url:        urls[i].Url,
			PreviewUrl: urls[i].PreviewUrl,
		}
	}

	dto := make([]dtos.FileDto, 0)

	for _, file := range newFiles {
		dto = append(dto, dtos.ToFileDto(file))
	}

	pagination.Total = searchResult.TotalHits
	pagination.Page = searchResult.Page
	pagination.Limit = searchResult.Limit

	return dto, pagination, err

}

func (s *FileService) DeleteFile(id uuid.UUID, userId uuid.UUID) error {
	file, err := s.fileRepo.GetFile(id)

	if err != nil {
		return err
	}

	folder, err := s.folderRepo.GetFolder(*file.FolderId)

	if err != nil {
		return err
	}

	err = s.addLogs("delete", userId, *file.FolderId, id, "", "", &folder.Name, nil)

	if err != nil {
		return err
	}

	_, err = s.fileRepo.DeleteFileDocument([]string{id.String()})

	if err != nil {
		return err
	}
	return s.fileRepo.DeleteFile(id)
}

func (s *FileService) UploadFileNoFolder(
	file multipart.File,
	fileName string,
	fileType string,
	userId uuid.UUID,
) (string, error) {
	casePermissions, err := s.casePermissionRepo.GetCasePermissionsByUserId(userId)

	if err != nil {
		return "", err
	}

	caseIds := make([]uuid.UUID, 0)
	cases := make([]models.Case, 0)

	for _, casePermission := range casePermissions {
		cases = append(cases, *casePermission.Case)
		caseIds = append(caseIds, casePermission.CaseId)
	}

	fileInDb, _ := s.fileRepo.GetFileByCaseId(fileName, caseIds)

	if len(fileInDb) > 0 {
		return s.uploadReplaceFile(file, fileType, fileInDb[0], userId)
	}

	mimeTypeString := convertMimeTypeToString(fileType)

	if mimeTypeString != "image" && mimeTypeString != "pdf" {
		return "", errors.New("there is no old file and we have no clue to identify")
	}
	link, err := s.uploadNewFileNoFolder(file, fileName, fileType, userId, cases)

	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *FileService) UploadFile(
	file multipart.File,
	fileName string,
	fileType string,
	folderId uuid.UUID,
	userId uuid.UUID,
) (string, error) {
	fileInDb, _ := s.getFileByNameInFolderId(fileName, folderId)

	if (*fileInDb).Name != "" {
		return s.uploadReplaceFile(file, fileType, *fileInDb, userId)
	}

	return s.uploadNewFile(file, fileName, fileType, folderId, userId)
}

func (s *FileService) uploadNewFile(
	file multipart.File,
	fileName string,
	fileType string,
	folderId uuid.UUID,
	userId uuid.UUID,
) (string, error) {

	caseId, err := s.checkPermissionAndGetCaseId(userId, folderId)
	mimeTypeToString := convertMimeTypeToString(fileType)

	if err != nil {
		return "", err
	}

	//needConvert := checkNeedConvert(fileType)

	version, versionPreview := "", ""

	cloudName := generateUniqueName(fileName)

	previewCloudName := cloudName

	s.logger.Info(cloudName, previewCloudName)

	version, err = s.storageService.UploadFile(file, cloudName)

	versionPreview = version

	if err != nil {
		s.logger.Error("Error: ", err)
		return "", err
	}

	fileT, err := s.fileTypeRepo.GetFileTypeByName(mimeTypeToString)

	if err != nil || fileT.ID == uuid.Nil {
		fileT, err = s.fileTypeRepo.GetEtcFileType()
	}

	tag, err := s.tagRepo.GetTagByNames([]string{mimeTypeToString})

	if err != nil {
		tagEtc, err := s.tagRepo.GetEtcTag()
		if err != nil {
			return "", err
		}
		tag = make([]models.Tag, 0)
		tag = append(tag, *tagEtc)
	}

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

	_, err = s.addMeili(fileRes, mimeTypeToString)

	if mimeTypeToString == "image" {
		go s.ocrImage(fileRes.ID, previewCloudName, tag, []models.Case{})
	} else if mimeTypeToString == "pdf" {
		go s.ocrPdf(fileRes.ID, previewCloudName, tag, []models.Case{})
	}

	if err != nil {
		return "", err
	}

	err = s.addLogs("upload", userId, *fileRes.FolderId, fileRes.ID, version, versionPreview, nil, nil)

	err = s.updateFolderTagByFolderId(folderId)

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return fileRes.ID.String(), nil
}

func (s *FileService) uploadReplaceFile(
	file multipart.File,
	fileType string,
	fileDb models.File,
	userId uuid.UUID,
) (string, error) {

	//needConvert := checkNeedConvert(fileType)

	version, versionPreview := "", ""

	cloudName := fileDb.CloudName

	version, err := s.storageService.UploadFile(file, cloudName)

	versionPreview = version

	if err != nil {
		return "", err
	}

	modelFile := models.File{}

	typeString := convertMimeTypeToString(fileType)

	err = s.fileRepo.UpdateFile(fileDb.ID, modelFile)
	s.logger.Info("TYPE: ", fileType)
	if typeString == "image" {
		go s.ocrImage(fileDb.ID, fileDb.CloudName, fileDb.Tags, []models.Case{})
	} else if typeString == "pdf" {
		go s.ocrPdf(fileDb.ID, fileDb.CloudName, fileDb.Tags, []models.Case{})
	}

	if err != nil {
		return "", err
	}

	err = s.addLogs("update", userId, *fileDb.FolderId, fileDb.ID, version, versionPreview, nil, nil)

	return "updated", nil
}

func (s *FileService) getFileByNameInFolderId(name string, folderId uuid.UUID) (*models.File, error) {
	return s.fileRepo.GetFileByName(name, folderId)
}

func generateUniqueName(fileName string) string {
	id := uuid.New().String()
	timestamp := time.Now().UnixNano()
	fileType := fileName[strings.LastIndex(fileName, ".")+1:]

	return fmt.Sprintf("%s-%d.%s", id, timestamp, fileType)
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

	s.logger.Info(fileViews)

	for i, fileView := range fileViews {
		files[i] = *fileView.File
	}

	fileUrls, err := s.storageService.GetSignedFileUrls(files)

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

func (s *FileService) addLogs(actionName string, userId, folderId, fileId uuid.UUID, version, previewVersion string, from, to *string) error {
	action, err := s.actionRepo.GetActionByName(actionName)

	if err != nil {
		return err
	}

	actionLog := models.ActionLog{
		FolderId:             folderId,
		FileId:               &fileId,
		UserId:               userId,
		ActionId:             action.ID,
		FileVersionId:        &version,
		FilePreviewVersionId: &previewVersion,
		From:                 from,
		To:                   to,
	}

	if err != nil {
		return err
	}

	_, err = s.actionLogRepo.CreateActionLog(actionLog)

	return err
}

func (s *FileService) addMeili(file models.File, mimeTypeToString string) (string, error) {

	tagString := make([]string, 0)
	for _, t := range file.Tags {
		tagString = append(tagString, t.Name)
	}
	folderString := make([]string, 0)

	if file.FolderId != nil {
		folderTilRoot, err := s.folderRepo.GetFromRootToFolder(*file.FolderId)
		if err != nil {
			return "", err
		}
		for _, f := range folderTilRoot {
			folderString = append(folderString, f.ID.String())
		}

	}

	caseString := ""

	if file.CaseId != nil {
		caseString = file.CaseId.String()
	}

	modelMeili := models.MeiliFile{
		Id:        file.ID.String(),
		Name:      file.Name,
		Type:      mimeTypeToString,
		Tags:      tagString,
		FolderIds: folderString,
		CaseId:    caseString,
	}

	res, err := s.fileRepo.CreateFileDocument(modelMeili)

	if err != nil {
		return "", err
	}

	s.logger.Info("RES ->", *res)

	return "", nil
}

func (s *FileService) updateFolderAndCaseInFile(id uuid.UUID, cases []models.Case, ocrData string) error {
	if len(cases) > 1 {
		s.logger.Info("More than one case found")
		var foundCase models.Case
		for _, c := range cases {
			titleSplit := strings.Split(c.Title, " ")
			for _, t := range titleSplit {
				if t != "" && strings.Contains(ocrData, t) {
					s.logger.Info("Found from [TITLE]->", t)
					foundCase = c
					break
				}
			}
			if foundCase.ID != uuid.Nil {
				break
			}
			if strings.Contains(ocrData, c.Title) {
				s.logger.Info("Found from ->", c.Title)
				foundCase = c
				break
			} else if c.RedCaseNumber != "" && strings.Contains(ocrData, c.RedCaseNumber) {
				s.logger.Info("Found from ->", c.RedCaseNumber)
				foundCase = c
				break
			} else if c.BlackCaseNumber != "" && strings.Contains(ocrData, c.BlackCaseNumber) {
				s.logger.Info("Found from ->", c.BlackCaseNumber)
				foundCase = c
				break
			}
		}
		if foundCase.ID != uuid.Nil {
			s.logger.Info("Found case ->", foundCase)

			//mCase, err := s.fileRepo.GetFileByName(name, foundCase.Folders[0].ID)

			var file models.File
			file.CaseId = &foundCase.ID
			file.FolderId = &foundCase.Folders[0].ID
			err := s.fileRepo.UpdateFile(id, file)

			if err != nil {
				return err
			}

			var fileMeili models.MeiliFile
			fileMeili.CaseId = foundCase.ID.String()
			fileMeili.Id = id.String()
			fileMeili.FolderIds = make([]string, 0)
			fileMeili.FolderIds = append(fileMeili.FolderIds, foundCase.Folders[0].ID.String())

			_, err = s.fileRepo.UpdateFileDocument(fileMeili)

		}
	}
	return nil
}

func (s *FileService) actionAfterOcr(id uuid.UUID, name string, tags []models.Tag, cases []models.Case, ocrData string) error {
	if ocrData == "" {
		return nil
	}

	go func() {
		err := s.updateFolderAndCaseInFile(id, cases, ocrData)

		if err != nil {
			s.logger.Error("Error while updating folder and case in file ->", err)
		}

		newTags := make([]models.Tag, len(tags))

		if strings.Contains(ocrData, "ทะเบียนบ้าน") {
			// Do nothing , because in home registration
			// There's id card number
			// So we need to detect it before to prevent wrong tagging
		} else if strings.Contains(ocrData, "บัตรประจำตัวประชาชน") {
			tag, _ := s.tagRepo.GetTagByName("idCard")
			if !utils.ContainsTag(tags, *tag) {
				newTags = append(tags, *tag)
			}
		}

		tagStrings := make([]string, 0)
		for _, t := range newTags {
			tagStrings = append(tagStrings, t.Name)
		}

		modelMeili := models.MeiliFile{
			Id:      id.String(),
			Content: ocrData,
			Tags:    tagStrings,
		}

		_, err = s.fileRepo.UpdateFileDocument(modelMeili)

		if err != nil {
			s.logger.Error(err)
		}

		if len(newTags) != len(tags) {
			err = s.fileRepo.UpdateTags(id, newTags)

			if err != nil {
				s.logger.Error(err)
			}

			file, err := s.fileRepo.GetFile(id)

			if err != nil {
				s.logger.Error(err)
			}

			err = s.updateFolderTagByFolderId(*file.FolderId)

			if err != nil {
				s.logger.Error(err)
			}
		}
	}()

	return nil
}

func (s *FileService) ocrImage(id uuid.UUID, name string, tags []models.Tag, cases []models.Case) {
	ocrData, err := s.ocrService.GetTextFromImageName(name)
	if err != nil {
		s.logger.Error(err)
	}

	err = s.actionAfterOcr(id, name, tags, cases, ocrData)

	if err != nil {
		s.logger.Error(err)
	}
}

func (s *FileService) ocrPdf(id uuid.UUID, name string, tags []models.Tag, cases []models.Case) {
	ocrData, err := s.ocrService.GetTextFromPdfUrl(name)

	if err != nil {
		s.logger.Error(err)
	}

	err = s.actionAfterOcr(id, name, tags, cases, ocrData)
}

func (s *FileService) uploadNewFileNoFolder(
	file multipart.File,
	fileName string,
	fileType string,
	userId uuid.UUID,
	cases []models.Case,
) (string, error) {

	mimeTypeToString := convertMimeTypeToString(fileType)

	cloudName := generateUniqueName(fileName)

	previewCloudName := cloudName

	version, err := s.storageService.UploadFile(file, cloudName)

	previewVersion := version

	if err != nil {
		return "", err
	}

	fileT, err := s.fileTypeRepo.GetFileTypeByName(mimeTypeToString)

	if err != nil {
		fileT, err = s.fileTypeRepo.GetEtcFileType()
	}

	tag, err := s.tagRepo.GetTagByNames([]string{mimeTypeToString})

	if err != nil {
		tagEtc, err := s.tagRepo.GetEtcTag()
		if err != nil {
			return "", err
		}
		tag = make([]models.Tag, 0)
		tag = append(tag, *tagEtc)
	}

	modelFile := models.File{
		Name:             fileName,
		TypeId:           &fileT.ID,
		CloudName:        cloudName,
		PreviewCloudName: previewCloudName,
		Tags:             tag,
	}

	fileRes, err := s.fileRepo.CreateFile(modelFile)

	_, err = s.addMeili(fileRes, mimeTypeToString)

	if mimeTypeToString == "image" {
		go s.ocrImage(fileRes.ID, previewCloudName, tag, cases)
	} else if mimeTypeToString == "pdf" {
		go s.ocrPdf(fileRes.ID, previewCloudName, tag, cases)
	}

	//if err != nil {
	//	return "", err
	//}

	err = s.addLogs("upload", userId, *fileRes.FolderId, fileRes.ID, version, previewVersion, nil, nil)

	//err = s.updateFolderTagByFolderId(folderId)

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return fileRes.ID.String(), nil
}
