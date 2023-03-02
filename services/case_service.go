package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type CaseService struct {
	logger   *libs.Logger
	caseRepo repositories.CaseRepository
}

func NewCaseService(logger *libs.Logger, r repositories.CaseRepository) CaseService {
	return CaseService{logger: logger, caseRepo: r}
}

func (s *CaseService) GetCases() (cases []models.Case, err error) {
	return s.caseRepo.GetCases()
}

func (s *CaseService) GetCase(id uint) (mCase *models.Case, err error) {
	return s.caseRepo.GetCase(id)
}

func (s *CaseService) CreateCase(mCase models.Case) (models.Case, error) {
	return s.caseRepo.CreateCase(mCase)
}

func (s *CaseService) UpdateCase(id uint, mCase models.Case) error {
	return s.caseRepo.UpdateCase(id, mCase)
}

func (s *CaseService) DeleteCase(id uint) error {
	return s.caseRepo.DeleteCase(id)
}
