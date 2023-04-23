package dtos

import (
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type FolderDto struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	SubFolders []FolderDto `json:"subFolders"`
	Files      []FileDto   `json:"files"`
	CreatedAt  string      `json:"createdAt"`
	UpdatedAt  string      `json:"updatedAt"`
	Tags       []TagDto    `json:"tags"`
}

func ToFolderDtos(folders []models.Folder) []FolderDto {
	var dtos []FolderDto

	for _, folder := range folders {
		dtos = append(dtos, ToFolderDto(folder))
	}

	return dtos
}

func ToFolderDto(folder models.Folder) FolderDto {
	return FolderDto{
		Id:         folder.ID.String(),
		Name:       folder.Name,
		Files:      ToFileDtos(folder.Files),
		CreatedAt:  folder.CreatedAt.String(),
		UpdatedAt:  folder.UpdatedAt.String(),
		SubFolders: ToFolderDtos(folder.SubFolders),
		Tags:       ToTagDtos(folder.Tags),
	}
}

type CreateFolderDto struct {
	Name           string `json:"name" binding:"required"`
	ParentFolderId string `json:"parentFolderId" binding:"required"`
}

type UpdateFolderDto struct {
	Name           string `json:"name"`
	ParentFolderId string `json:"parentFolderId"`
}

func (dto CreateFolderDto) ToModel(caseId uuid.UUID) models.Folder {

	parentFolderId, err := uuid.Parse(dto.ParentFolderId)

	if err != nil {
		return models.Folder{
			Name:           dto.Name,
			ParentFolderId: nil,
			IsArchive:      false,
			CaseId:         &caseId,
		}
	}

	return models.Folder{
		Name:           dto.Name,
		ParentFolderId: &parentFolderId,
		IsArchive:      false,
		CaseId:         &caseId,
	}
}

func (dto UpdateFolderDto) ToModel() models.Folder {

	parentFolderId, err := uuid.Parse(dto.ParentFolderId)

	if err != nil {
		return models.Folder{
			Name: dto.Name,
		}
	}

	return models.Folder{
		Name:           dto.Name,
		ParentFolderId: &parentFolderId,
	}
}
