package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type CaseService struct {
	logger     *libs.Logger
	caseRepo   repositories.CaseRepository
	folderRepo repositories.FolderRepository
}

func NewCaseService(logger *libs.Logger, r repositories.CaseRepository, f repositories.FolderRepository) CaseService {
	return CaseService{logger: logger, caseRepo: r, folderRepo: f}
}

func (s *CaseService) GetCases() (cases []models.Case, err error) {
	return s.caseRepo.GetCases()
}

func (s *CaseService) GetCase(id uint) (mCase *models.Case, err error) {
	return s.caseRepo.GetCase(id)
}

func (s *CaseService) CreateCase(dto dtos.CreateCaseDto) (models.Case, error) {

	mCase := dto.ToCase(models.Folder{
		Name:      dto.Name,
		IsArchive: false,
	})

	return s.caseRepo.CreateCase(mCase)
}

func (s *CaseService) UpdateCase(id uint, mCase models.Case) error {
	return s.caseRepo.UpdateCase(id, mCase)
}

func (s *CaseService) DeleteCase(id uint) error {
	return s.caseRepo.DeleteCase(id)
}
