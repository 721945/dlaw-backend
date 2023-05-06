package dtos

import "github.com/721945/dlaw-backend/models"

type AddMemberDto struct {
	UserIds    []string `json:"userIds" binding:"required"`
	Permission string   `json:"permission" binding:"required"`
}

type UpdateMemberDto struct {
	Permission string `json:"permission" binding:"required"`
}

type MemberDto struct {
	Firstname  string `json:"firstName"`
	Lastname   string `json:"lastName"`
	Permission string `json:"permission"`
	Email      string `json:"email"`
}

func ToMemberDto(permission models.CasePermission) MemberDto {
	return MemberDto{
		Firstname:  permission.User.Firstname,
		Lastname:   permission.User.Lastname,
		Permission: permission.Permission.Name,
		Email:      permission.User.Email,
	}
}
