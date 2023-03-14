package dtos

import "github.com/721945/dlaw-backend/models"

type PermissionDto struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type CreatePermissionDto struct {
	Name string `json:"name"`
}

type UpdatePermissionDto struct {
	Name string `json:"name"`
}

func (p CreatePermissionDto) ToModel() models.Permission {
	return models.Permission{
		Name: p.Name,
	}
}

func (p UpdatePermissionDto) ToModel() models.Permission {
	return models.Permission{
		Name: p.Name,
	}
}

func ToPermissionDto(permission models.Permission) PermissionDto {
	return PermissionDto{
		Id:   permission.ID,
		Name: permission.Name,
	}
}
