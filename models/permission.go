package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name            string
	CasePermissions []CasePermission    `gorm:"foreignKey:PermissionId"`
	PermissionLogs  []CasePermissionLog `gorm:"foreignKey:PermissionId"`
}
