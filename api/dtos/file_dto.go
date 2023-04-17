package dtos

import (
	"github.com/721945/dlaw-backend/models"
)

type CreateFileDto struct {
	Name string `json:"name" binding:"required"`
	Url  string `json:"url" binding:"required"`
}

type FileDto struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	PreviewUrl string `json:"previewUrl"`
}

func ToFileDto(file models.File) FileDto {
	return FileDto{
		Id:   file.ID.String(),
		Name: file.Name,
		//Url:  file.Url.Url,
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
