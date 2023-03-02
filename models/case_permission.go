package models

import "gorm.io/gorm"

type CasePermission struct {
	gorm.Model
	UserId       uint
	CaseId       uint
	PermissionId uint
}
