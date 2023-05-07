package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
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
	caseUsedRepo       repositories.CaseUsedLogRepository
	permissionLogRepo  repositories.CasePermissionLogRepository
	userRepo           repositories.UserRepository
	storage            google_storage.GoogleStorage
}

func NewCaseService(
	logger *libs.Logger,
	r repositories.CaseRepository,
	f repositories.FolderRepository,
	p repositories.PermissionRepository,
	cp repositories.CasePermissionRepository,
	pl repositories.CasePermissionLogRepository,
	u repositories.UserRepository,
	cup repositories.CaseUsedLogRepository,
	storage google_storage.GoogleStorage,
) CaseService {
	return CaseService{
		logger:             logger,
		caseRepo:           r,
		folderRepo:         f,
		permissionRepo:     p,
		casePermissionRepo: cp,
		permissionLogRepo:  pl,
		userRepo:           u,
		caseUsedRepo:       cup,
		storage:            storage,
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

	for i, mCase := range cases {
		casesDto[i] = dtos.ToCaseDto(mCase)
	}

	return casesDto, nil
}

func (s CaseService) checkPermission(caseId uuid.UUID, userId uuid.UUID) (bool, error) {
	permissionCase, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(userId, caseId)

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

	for i, mCase := range cases {
		casesDto[i] = dtos.ToCaseDto(mCase)
	}

	return casesDto, nil
}

func (s CaseService) GetPublicCases() (casesDto []dtos.CasePublicDto, err error) {
	cases, err := s.caseRepo.GetCasesWhichFileIsPublic()

	if err != nil {
		return casesDto, err
	}

	//var wg sync.WaitGroup
	for _, myCase := range cases {
		//wg.Add(1)
		//go func(myCase models.Case) {
		urls, err := s.storage.GetSignedFileUrls(myCase.Files)

		s.logger.Info(urls)
		if err != nil {
			s.logger.Error(err)
			continue
		}

		//newFiles := make([]models.File, len(myCase.Files))

		for i, _ := range myCase.Files {
			//newFiles[i] = file
			myCase.Files[i].Url = &models.FileUrl{
				Url:        urls[i].Url,
				PreviewUrl: urls[i].PreviewUrl,
			}
		}

		//myCase.Files = newFiles

		//}(myCase)
	}

	//wg.Wait()

	casesDto = make([]dtos.CasePublicDto, len(cases))
	//
	for i, mCase := range cases {
		casesDto[i] = dtos.ToCasePublicDto(mCase)
	}

	return casesDto, nil
}

func (s CaseService) GetFrequencyCases(userId uuid.UUID) (casesDto []dtos.CaseDto, err error) {
	permissionCases, err := s.casePermissionRepo.GetCasePermissionsByUserId(userId)

	if err != nil {
		return casesDto, err
	}

	caseIds := make([]uuid.UUID, len(permissionCases))

	for i, permissionCase := range permissionCases {
		caseIds[i] = permissionCase.CaseId
	}

	casesUsed, err := s.caseUsedRepo.GetCaseUsedLogWithCaseByCaseIdsAndUserId(caseIds, userId)

	//cases, err := s.caseRepo.GetCasesSortedByFrequency(userId)

	if err != nil {
		return casesDto, err
	}

	casesDto = make([]dtos.CaseDto, len(casesUsed))

	for i, mCase := range casesUsed {
		casesDto[i] = dtos.ToSimpleCaseDto(*mCase.Case)
	}

	return casesDto, nil
}

func (s CaseService) GetFolders(caseId uuid.UUID) (dto []dtos.SimpleFolderDto, err error) {
	folders, err := s.folderRepo.GetFoldersByCaseId(caseId)

	if err != nil {
		return dto, err
	}

	dto = make([]dtos.SimpleFolderDto, len(folders))

	for i, folder := range folders {
		dto[i] = dtos.SimpleFolderDto{
			Id:   folder.ID.String(),
			Name: folder.Name,
		}
	}

	return dto, nil
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

func (s CaseService) AddMember(caseId uuid.UUID, dto dtos.AddMemberDto) (ids []string, err error) {
	permission, err := s.permissionRepo.GetPermissionByName(dto.Permission)

	if err != nil {
		return []string{}, err
	}

	userIds := make([]uuid.UUID, len(dto.UserIds))

	for i, userId := range dto.UserIds {
		userIds[i], err = uuid.Parse(userId)
	}

	if err != nil {
		return []string{}, err
	}

	//user, err := s.userRepo.GetUser(userId)
	//
	//permissionCase := models.CasePermission{
	//	CaseId:       caseId,
	//	PermissionId: permission.ID,
	//	UserId:       user.ID,
	//}
	//
	//casePermission, err := s.casePermissionRepo.CreateCasePermission(permissionCase)
	resCh := make(chan string)

	// Start a goroutine for each user to create a case permission concurrently
	for _, uid := range userIds {
		go func(userId uuid.UUID) {
			permissionItem, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(userId, caseId)

			if err != nil {
				resCh <- ""
			}

			if permissionItem != nil {
				permissionCase := models.CasePermission{
					CaseId:       caseId,
					PermissionId: permission.ID,
					UserId:       userId,
				}
				cpm, err := s.casePermissionRepo.CreateCasePermission(permissionCase)
				if err != nil {
					resCh <- ""
				} else {
					idString := cpm.ID.String()
					resCh <- idString
				}
			} else {
				permissionCase := models.CasePermission{
					CaseId:       caseId,
					PermissionId: permission.ID,
					UserId:       userId,
				}
				err := s.casePermissionRepo.UpdateCasePermission(permissionItem.ID, permissionCase)

				if err != nil {
					resCh <- ""
				} else {
					idString := permission.ID.String()
					resCh <- idString
				}
			}
		}(uid)
	}

	// Collect the results of the CreateCasePermission calls
	casePermissionIds := make([]string, len(userIds))
	for i := 0; i < len(userIds); i++ {
		casePermissionId := <-resCh
		casePermissionIds[i] = casePermissionId
	}

	return casePermissionIds, err
}

func (s CaseService) RemoveMember(caseId uuid.UUID, userId uuid.UUID) (err error) {
	permission, err := s.casePermissionRepo.GetCasePermissionsByUserIdAndCaseId(userId, caseId)

	if err != nil {
		return err
	}

	err = s.casePermissionRepo.DeleteCasePermission(permission.ID)

	return err
}
