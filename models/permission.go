package models

type Permission struct {
	Base
	Name            string
	CasePermissions []CasePermission    `gorm:"foreignKey:PermissionId"`
	PermissionLogs  []CasePermissionLog `gorm:"foreignKey:PermissionId"`
}
