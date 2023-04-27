package dtos

import (
	"github.com/721945/dlaw-backend/models"
	"time"
)

type CaseDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CaseDetailDto struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Tags      []TagDto  `json:"tags"`
	Owner     UserDto   `json:"owner"`
	ShareWith []UserDto `json:"shareWith"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FolderId  string    `json:"folderId"`
}

type OwnerDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateCaseDto struct {
	RedCaseNumber   string  `json:"redCaseNumber"`
	BlackCaseNumber string  `json:"blackCaseNumber" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Email           *string `json:"email" binding:"required"`
	CaseTitle       *string `json:"caseTitle" binding:"required"`
	CaseContent     *string `json:"caseContent" binding:"required"`
}

type UpdateCaseDto struct {
	RedCaseNumber   string  `json:"redCaseNumber" `
	BlackCaseNumber string  `json:"blackCaseNumber"`
	Name            string  `json:"name"`
	Email           *string `json:"email"`
	CaseTitle       *string `json:"caseTitle"`
	CaseContent     *string `json:"caseContent"`
}

func (dto CreateCaseDto) ToCase(folder models.Folder) models.Case {
	return models.Case{
		RedCaseNumber:   dto.RedCaseNumber,
		BlackCaseNumber: dto.BlackCaseNumber,
		Title:           dto.Name,
		CaseTitle:       dto.CaseTitle,
		CaseDetail:      dto.CaseContent,
		Folders:         []models.Folder{folder},
		Email:           dto.Email,
	}
}

func (dto UpdateCaseDto) ToModel() models.Case {
	return models.Case{
		RedCaseNumber:   dto.RedCaseNumber,
		BlackCaseNumber: dto.BlackCaseNumber,
		Title:           dto.Name,
		CaseTitle:       dto.CaseTitle,
		CaseDetail:      dto.CaseContent,
		Email:           dto.Email,
	}
}

func ToCaseDto(mCase models.Case) CaseDetailDto {
	permissions := mCase.CasePermissions

	users := make([]UserDto, len(permissions))

	var owner UserDto

	for i, permission := range permissions {
		if permission.Permission.Name == "owner" {
			owner = UserDto{
				ID:        permission.User.ID.String(),
				FirstName: permission.User.Firstname,
				LastName:  permission.User.Lastname,
				Email:     permission.User.Email,
			}
		}
		users[i] = UserDto{
			ID:        permission.User.ID.String(),
			FirstName: permission.User.Firstname,
			LastName:  permission.User.Lastname,
			Email:     permission.User.Email,
		}
	}

	return CaseDetailDto{
		Id:        mCase.ID.String(),
		Name:      mCase.Title,
		Tags:      ToTagDtos(mCase.Folders[0].Tags),
		Owner:     owner,
		ShareWith: users,
		CreatedAt: mCase.CreatedAt,
		UpdatedAt: mCase.UpdatedAt,
		FolderId:  mCase.Folders[0].ID.String(),
	}
}
