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

	folder, err := s.folderRepo.GetRootFolder(id)

	if err != nil {
		return err
	}

	mCase := dto.ToModel()

	if mCase.Title != "" {
		folder.Name = mCase.Title

		err = s.folderRepo.UpdateFolder(folder.ID, *folder)
	}

	err = s.caseRepo.UpdateCase(id, mCase)

	return err
}

func (s CaseService) DeleteCase(id uuid.UUID, userId uuid.UUID) error {
	can, err := s.checkPermission(id, userId)

	if err != nil {
		return err
	}

	if !can {
		return libs.ErrUnauthorized
	}

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

	cases, err := s.caseRepo.GetCasesByIds(caseIds, false)

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

func (s CaseService) ArchiveCase(id uuid.UUID, userId uuid.UUID) error {
	can, err := s.checkPermission(id, userId)

	if err != nil {
		return err
	}

	if !can {
		return libs.ErrUnauthorized
	}

	return s.caseRepo.ArchiveCase(id)
}
func (s CaseService) UnArchiveCase(id uuid.UUID, userId uuid.UUID) error {
	can, err := s.checkPermission(id, userId)

	if err != nil {
		return err
	}

	if !can {
		return libs.ErrUnauthorized
	}

	return s.caseRepo.ArchiveCase(id)
}

func (s CaseService) GetArchivedCases(userId uuid.UUID) (casesDto []dtos.CaseDetailDto, err error) {
	permissionCases, err := s.casePermissionRepo.GetCasePermissionsByUserId(userId)

	if err != nil {
		return casesDto, err
	}

	caseIds := make([]uuid.UUID, len(permissionCases))

	for i, permissionCase := range permissionCases {
		caseIds[i] = permissionCase.CaseId
	}

	cases, err := s.caseRepo.GetCasesByIds(caseIds, true)

	if err != nil {
		return casesDto, err
	}

	casesDto = make([]dtos.CaseDetailDto, len(cases))

	for _, mCase := range cases {
		casesDto = append(casesDto, dtos.ToCaseDto(mCase))
	}

	return casesDto, nil
}

func (s CaseService) GetFrequencyCases(userId uuid.UUID) (casesDto []dtos.CaseDetailDto, err error) {
	permissionCases, err := s.casePermissionRepo.GetCasePermissionsByUserId(userId)

	if err != nil {
		return casesDto, err
	}

	caseIds := make([]uuid.UUID, len(permissionCases))

	for i, permissionCase := range permissionCases {
		caseIds[i] = permissionCase.CaseId
	}

	cases, err := s.caseRepo.GetCasesByIds(caseIds, false)

	if err != nil {
		return casesDto, err
	}

	casesDto = make([]dtos.CaseDetailDto, len(cases))

	for _, mCase := range cases {
		casesDto = append(casesDto, dtos.ToCaseDto(mCase))
	}

	return casesDto, nil
}

func (s CaseService) GetMembers(caseId uuid.UUID) (members []dtos.MemberDto, err error) {
	permissions, err := s.casePermissionRepo.GetCasePermissionsByCaseId(caseId)

	if err != nil {
		return members, err
	}

	members = make([]dtos.MemberDto, len(permissions))

	for i, permission := range permissions {
		members[i] = dtos.ToMemberDto(permission)
	}

	return members, nil
}

func (s CaseService) UpdateMember(caseId uuid.UUID, userId uuid.UUID, dto dtos.UpdateMemberDto) (err error) {
	permission, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(caseId, userId)

	if err != nil {
		return err
	}

	permissionUpdate, err := s.permissionRepo.GetPermissionByName(dto.Permission)

	if err != nil {
		return err
	}

	err = s.casePermissionRepo.UpdateCasePermission(permission.ID, models.CasePermission{PermissionId: permissionUpdate.ID})

	return err
}

func (s CaseService) DeleteMember(caseId uuid.UUID, userId uuid.UUID) (err error) {
	permission, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(caseId, userId)

	if err != nil {
		return err
	}

	err = s.casePermissionRepo.DeleteCasePermission(permission.ID)

	return err
}
