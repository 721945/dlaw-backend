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
	Owner     OwnerDto  `json:"owner"`
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
	BlackCaseNumber string  `json:"blackCaseNumber"`
	Name            string  `json:"name"`
	Email           *string `json:"email"`
	CaseTitle       *string `json:"caseTitle"`
	CaseContent     *string `json:"caseContent"`
}

type UpdateCaseDto struct {
	RedCaseNumber   string  `json:"redCaseNumber"`
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
	return CaseDetailDto{
		Id:   mCase.ID.String(),
		Name: mCase.Title,
		Tags: ToTagDtos(mCase.Folders[0].Tags),
		Owner: OwnerDto{
			FirstName: "John",
			LastName:  "Doe",
		},
		ShareWith: []UserDto{},
		CreatedAt: mCase.CreatedAt,
		UpdatedAt: mCase.UpdatedAt,
		FolderId:  mCase.Folders[0].ID.String(),
	}
}
