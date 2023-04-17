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

func (s *PermissionService) GetPermissions() (permissions []models.Permission, err error) {
	return s.permissionRepo.GetPermissions()

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

func (s *PermissionService) CreatePermission(permission models.Permission) (models.Permission, error) {
	return s.permissionRepo.CreatePermission(permission)
}

func (s *PermissionService) UpdatePermission(id uuid.UUID, permission models.Permission) error {
	return s.permissionRepo.UpdatePermission(id, permission)
}

func (s *PermissionService) DeletePermission(id uuid.UUID) error {
	return s.permissionRepo.DeletePermission(id)
}
