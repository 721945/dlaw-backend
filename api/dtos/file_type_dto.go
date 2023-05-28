package dtos

import (
	"github.com/721945/dlaw-backend/models"
)

type FileTypeDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateFileTypeDto struct {
	Name string `json:"name" binding:"required"`
}

type CreateFileTypesDto struct {
	NameList []string `json:"nameList" binding:"required"`
}

type UpdateFileTypeDto struct {
	Name string `json:"name" binding:"required"`
}

func (dto CreateFileTypeDto) ToFileType() models.FileType {
	return models.FileType{
		Name: dto.Name,
	}
}

func (dto UpdateFileTypeDto) ToFileType() models.FileType {
	return models.FileType{
		Name: dto.Name,
	}
}

func ToFileType(fileType *models.FileType) *FileTypeDto {
	return &FileTypeDto{
		ID:   fileType.ID.String(),
		Name: fileType.Name,
	}
}

func ToFileTypes(fileTypes []models.FileType) []FileTypeDto {
	var fileTypeDtos = make([]FileTypeDto, 0)

	for _, fileType := range fileTypes {
		fileTypeDtos = append(fileTypeDtos, *ToFileType(&fileType))
	}

	return fileTypeDtos
}
