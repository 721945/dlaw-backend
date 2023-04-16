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

func (dto CreateCaseDto) ToCase(folder models.Folder) models.Case {
	return models.Case{
		RedCaseNumber:   dto.RedCaseNumber,
		BlackCaseNumber: dto.BlackCaseNumber,
		Title:           dto.Name,
		CaseTitle:       dto.CaseTitle,
		CaseDetail:      dto.CaseContent,
		Folder:          folder,
	}
}
