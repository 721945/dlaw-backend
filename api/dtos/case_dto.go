package dtos

import "github.com/721945/dlaw-backend/models"

type CaseDto struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCaseDto struct {
	CaseNumber  string  `json:"caseNumber"`
	Name        string  `json:"name"`
	Email       *string `json:"email"`
	CaseTitle   *string `json:"caseTitle"`
	CaseContent *string `json:"caseContent"`
}

func (dto CreateCaseDto) ToCase(folder models.Folder) models.Case {
	return models.Case{
		CaseNumber: dto.CaseNumber,
		Title:      dto.Name,
		CaseTitle:  dto.CaseTitle,
		CaseDetail: dto.CaseContent,
		Folder:     folder,
	}
}
