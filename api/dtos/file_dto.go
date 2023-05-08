package dtos

import (
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type CreateFileDto struct {
	Name string `json:"name" binding:"required"`
}

type FileDto struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	PreviewUrl string   `json:"previewUrl"`
	Tags       []TagDto `json:"tags,omitempty"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
	Type       string   `json:"type,omitempty"`
	Size       string   `json:"size,omitempty"`
	IsPublic   bool     `json:"isPublic,omitempty"`
	IsShare    bool     `json:"isShare,omitempty"`
}

type FilePublicDto struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	PreviewUrl string `json:"previewUrl"`
	Type       string `json:"type,omitempty"`
}

func ToFilePublicDto(file models.File) FilePublicDto {
	var url, previewUrl, fileType string
	if file.Url != nil {
		url = file.Url.Url
		previewUrl = file.Url.PreviewUrl
	}
	if file.FileType != nil {
		fileType = file.FileType.Name
	}

	return FilePublicDto{
		Id:         file.ID.String(),
		Name:       file.Name,
		Url:        url,
		PreviewUrl: previewUrl,
		Type:       fileType,
	}
}

func ToFilePublicDtos(files []models.File) []FilePublicDto {
	filesDto := make([]FilePublicDto, len(files))

	for i, file := range files {
		filesDto[i] = ToFilePublicDto(file)
	}

	return filesDto
}

func ToFileDto(file models.File) FileDto {
	var url, previewUrl, fileType string
	if file.Url != nil {
		url = file.Url.Url
		previewUrl = file.Url.PreviewUrl
	}
	if file.FileType != nil {
		fileType = file.FileType.Name
	}
	return FileDto{
		Id:         file.ID.String(),
		Name:       file.Name,
		Url:        url,
		PreviewUrl: previewUrl,
		Tags:       ToTagDtos(file.Tags),
		CreatedAt:  file.CreatedAt.String(),
		UpdatedAt:  file.UpdatedAt.String(),
		Type:       fileType,
	}
}

func ToFileWithSizeDto(file models.File, size int64) FileDto {
	var url, previewUrl, fileType string
	if file.Url != nil {
		url = file.Url.Url
		previewUrl = file.Url.PreviewUrl
	}
	if file.FileType != nil {
		fileType = file.FileType.Name
	}

	return FileDto{
		Id:         file.ID.String(),
		Name:       file.Name,
		Url:        url,
		PreviewUrl: previewUrl,
		Tags:       ToTagDtos(file.Tags),
		CreatedAt:  file.CreatedAt.String(),
		UpdatedAt:  file.UpdatedAt.String(),
		Type:       fileType,
		Size:       utils.FormatFileSize(size),
		IsPublic:   file.IsPublic,
		IsShare:    file.IsShared,
	}
}

func ToFileDtos(files []models.File) []FileDto {
	filesDto := make([]FileDto, len(files))

	for i, file := range files {
		filesDto[i] = ToFileDto(file)
	}

	return filesDto
}

func (c CreateFileDto) ToModel() models.File {
	return models.File{
		Name: c.Name,
	}
}

type UpdateFileDto struct {
	Name string `json:"name"`
}

func (u UpdateFileDto) ToModel() models.File {
	return models.File{
		Name: u.Name,
	}
}

type GetSingleFileDto struct {
	Amount int `json:"amount"`
}

type MoveFileDto struct {
	TargetFolderId string `json:"targetFolderId" binding:"required"`
}

func (m MoveFileDto) ToModel() *models.File {
	id, err := uuid.Parse(m.TargetFolderId)
	if err != nil {
		return nil
	}

	return &models.File{
		FolderId: &id,
	}
}
