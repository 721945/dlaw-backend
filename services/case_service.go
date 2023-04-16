package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type CaseService struct {
	logger             *libs.Logger
	caseRepo           repositories.CaseRepository
	folderRepo         repositories.FolderRepository
	permissionRepo     repositories.PermissionRepository
	casePermissionRepo repositories.CasePermissionRepository
	permissionLogRepo  repositories.CasePermissionLogRepository
}

func NewCaseService(logger *libs.Logger, r repositories.CaseRepository, f repositories.FolderRepository) CaseService {
	return CaseService{logger: logger, caseRepo: r, folderRepo: f}
}

func (s CaseService) GetCases() (cases []models.Case, err error) {
	return s.caseRepo.GetCases()
}

func (s CaseService) GetCase(id uuid.UUID) (mCase *models.Case, err error) {
	return s.caseRepo.GetCase(id)
}

func (s CaseService) CreateCase(dto dtos.CreateCaseDto, userId uuid.UUID) (models.Case, error) {

	mCase := dto.ToCase(models.Folder{
		Name:      dto.Name,
		IsArchive: false,
	})

	permission, err := s.permissionRepo.GetPermissionByName("owner")

	if err != nil {
		return mCase, err
	}

	mCase, err = s.caseRepo.CreateCase(mCase)

	if err != nil {
		return mCase, err
	}

	casePermission := models.CasePermission{
		UserId:       userId,
		CaseId:       mCase.ID,
		PermissionId: permission.ID,
	}

	_, err = s.casePermissionRepo.CreateCasePermission(casePermission)

	if err != nil {
		return mCase, err
	}

	return mCase, nil
}

func (s CaseService) UpdateCase(id uuid.UUID, mCase models.Case) error {
	return s.caseRepo.UpdateCase(id, mCase)
}

func (s CaseService) DeleteCase(id uuid.UUID) error {
	return s.caseRepo.DeleteCase(id)
}

func (s CaseService) GetOwnCases(id uuid.UUID) (cases []dtos.CaseDetailDto, err error) {

	permissionCases, err := s.casePermissionRepo.GetCasePermissionsByUserId(id)

	s.logger.Info(permissionCases)
	//
	//if err != nil {
	//	return cases, err
	//}

	//for _, permissionCase := range permissionCases {
	//	//mCase, err := s.caseRepo.GetCase(permissionCase.CaseId)
	//
	//	if err != nil {
	//		return cases, err
	//	}
	//
	//	//cases = append(cases, dtos.CaseDetailDto{
	//	//	ID:        mCase.ID,
	//	//	Name:      mCase.Name,
	//	//	IsArchive: mCase.IsArchive,
	//	//	IsPrivate: mCase.IsPrivate,
	//	//	IsPublic:  mCase.IsPublic,
	//	//	UserId:    mCase.UserId,
	//	//})
	//}

	return cases, nil
}
