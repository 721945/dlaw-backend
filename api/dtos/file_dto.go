package dtos

import (
	"github.com/721945/dlaw-backend/models"
)

type CreateFileDto struct {
	Name string `json:"name" binding:"required"`
	Url  string `json:"url" binding:"required"`
}

type FileDto struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	PreviewUrl string   `json:"previewUrl"`
	Tags       []TagDto `json:"tags,omitempty"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

func ToFileDto(file models.File) FileDto {
	var url, previewUrl string
	if file.Url != nil {
		url = file.Url.Url
		previewUrl = file.Url.PreviewUrl
	}
	return FileDto{
		Id:         file.ID.String(),
		Name:       file.Name,
		Url:        url,
		PreviewUrl: previewUrl,
		Tags:       ToTagDtos(file.Tags),
		CreatedAt:  file.CreatedAt.String(),
		UpdatedAt:  file.UpdatedAt.String(),
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
	Url  string `json:"url"`
}

func (u UpdateFileDto) ToModel() models.File {
	return models.File{
		Name: u.Name,
	}
}

type GetSingleFileDto struct {
	Amount int `json:"amount"`
}
