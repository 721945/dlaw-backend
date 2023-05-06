package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type PermissionService struct {
	logger         *libs.Logger
	permissionRepo repositories.PermissionRepository
}

func NewPermissionService(logger *libs.Logger, r repositories.PermissionRepository) PermissionService {
	return PermissionService{logger: logger, permissionRepo: r}

}

func (s *PermissionService) GetPermissions() (permissions []dtos.PermissionDto, err error) {
	permissionModels, err := s.permissionRepo.GetPermissions()

	permissions = make([]dtos.PermissionDto, len(permissionModels))

	for i, permission := range permissionModels {
		permissions[i] = *dtos.ToPermissionDto(&permission)
	}
	return permissions, err
}

func (s *PermissionService) GetPermission(id uuid.UUID) (permissionDto *dtos.PermissionDto, err error) {
	permission, err := s.permissionRepo.GetPermission(id)

	if err != nil {
		return permissionDto, err
	}

	permissionDto = dtos.ToPermissionDto(permission)

	return permissionDto, nil
}

func (s *PermissionService) GetPermissionByName(name string) (permissionDto *dtos.PermissionDto, err error) {
	permission, err := s.permissionRepo.GetPermissionByName(name)

	if err != nil {
		return permissionDto, err
	}

	permissionDto = dtos.ToPermissionDto(permission)

	return permissionDto, nil
}

func (s *PermissionService) CreatePermission(permission models.Permission) (string, error) {
	permission, err := s.permissionRepo.CreatePermission(permission)

	return permission.ID.String(), err
}

func (s *PermissionService) UpdatePermission(id uuid.UUID, permission models.Permission) error {
	return s.permissionRepo.UpdatePermission(id, permission)
}

func (s *PermissionService) DeletePermission(id uuid.UUID) error {
	return s.permissionRepo.DeletePermission(id)
}
