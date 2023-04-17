package models

import "github.com/google/uuid"

type CasePermission struct {
	Base
	UserId       uuid.UUID
	User         *User
	CaseId       uuid.UUID
	Case         *Case
	PermissionId uuid.UUID
	Permission   *Permission
}
