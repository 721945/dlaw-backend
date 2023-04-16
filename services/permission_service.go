package services

import (
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

func (s *PermissionService) GetPermission(id uuid.UUID) (permission *models.Permission, err error) {
	return s.permissionRepo.GetPermission(id)

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
