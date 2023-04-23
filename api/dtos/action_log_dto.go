package dtos

import "github.com/721945/dlaw-backend/models"

type ActionLogDto struct {
	Id        string            `json:"id"`
	FolderId  string            `json:"folderId"`
	Action    string            `json:"action"`
	User      UserDto           `json:"user"`
	CreatedAt string            `json:"createdAt"`
	File      *FileActionLogDto `json:"file,omitempty"`
}

type FileActionLogDto struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	FileUrl  string `json:"fileUrl"`
}

func ToActionLogDto(actionLog models.ActionLog, file *FileActionLogDto) ActionLogDto {
	return ActionLogDto{
		Id:        actionLog.ID.String(),
		FolderId:  actionLog.FolderId.String(),
		Action:    actionLog.Action.Name,
		User:      *ToUserDto(actionLog.User),
		CreatedAt: actionLog.CreatedAt.Format("2006-01-02 15:04:05"),
		File:      file,
	}
}

func ToActionLogDtoList(actionLogs []models.ActionLog) []ActionLogDto {
	var actionLogDtos []ActionLogDto
	for _, actionLog := range actionLogs {
		actionLogDtos = append(actionLogDtos, ToActionLogDto(actionLog, nil))
	}
	return actionLogDtos
}

func ToFileActionLogDto(file *models.File, url string) *FileActionLogDto {
	if file == nil {
		return nil
	}

	return &FileActionLogDto{
		Id:       file.ID.String(),
		Filename: file.Name,
		FileUrl:  url,
	}
}
