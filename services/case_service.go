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

func NewCaseService(
	logger *libs.Logger,
	r repositories.CaseRepository,
	f repositories.FolderRepository,
	p repositories.PermissionRepository,
	cp repositories.CasePermissionRepository,
	pl repositories.CasePermissionLogRepository,
) CaseService {
	return CaseService{
		logger:             logger,
		caseRepo:           r,
		folderRepo:         f,
		permissionRepo:     p,
		casePermissionRepo: cp,
		permissionLogRepo:  pl,
	}
}

func (s CaseService) GetCases() (cases []models.Case, err error) {
	return s.caseRepo.GetCases()
}

func (s CaseService) GetCase(id, userId uuid.UUID) (mCase *models.Case, err error) {
	return s.caseRepo.GetCase(id)
}

func (s CaseService) CreateCase(dto dtos.CreateCaseDto, userId uuid.UUID) (string, error) {

	mCase := dto.ToCase(models.Folder{
		Name:      dto.Name,
		IsArchive: false,
	})

	permission, err := s.permissionRepo.GetPermissionByName("owner")

	if err != nil {
		return "", err
	}

	mCase, err = s.caseRepo.CreateCase(mCase)

	if err != nil {
		return "", err

	}

	casePermission := models.CasePermission{
		UserId:       userId,
		CaseId:       mCase.ID,
		PermissionId: permission.ID,
	}

	s.logger.Info(casePermission)

	_, err = s.casePermissionRepo.CreateCasePermission(casePermission)

	if err != nil {
		return "", err
	}

	return mCase.ID.String(), nil
}

func (s CaseService) UpdateCase(id uuid.UUID, dto dtos.UpdateCaseDto, userId uuid.UUID) error {
	// Need to check permission before
	hasPermission, err := s.checkPermission(id, userId)

	if err != nil {
		return err
	}

	if !hasPermission {
		return libs.ErrUnauthorized
	}

	mCase := dto.ToModel()

	err = s.caseRepo.UpdateCase(id, mCase)

	//if err != nil {
	//	return err
	//}

	return err
}

func (s CaseService) DeleteCase(id uuid.UUID) error {
	return s.caseRepo.DeleteCase(id)
}

func (s CaseService) GetOwnCases(id uuid.UUID) (casesDto []dtos.CaseDetailDto, err error) {

	permissionCases, err := s.casePermissionRepo.GetCasePermissionsByUserId(id)

	if err != nil {
		return casesDto, err
	}

	caseIds := make([]uuid.UUID, len(permissionCases))

	for i, permissionCase := range permissionCases {
		caseIds[i] = permissionCase.CaseId
	}

	cases, err := s.caseRepo.GetCasesByIds(caseIds)

	if err != nil {
		return casesDto, err
	}

	casesDto = make([]dtos.CaseDetailDto, len(cases))

	for _, mCase := range cases {
		casesDto = append(casesDto, dtos.ToCaseDto(mCase))
	}

	return casesDto, nil
}

func (s CaseService) checkPermission(caseId uuid.UUID, userId uuid.UUID) (bool, error) {
	permissionCase, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(caseId, userId)

	if err != nil {
		return false, err
	}

	if permissionCase.Permission.Name == "owner" {
		return true, nil
	}

	return false, nil
}
