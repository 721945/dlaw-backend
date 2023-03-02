package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Firstname            string
	Lastname             string
	Email                string
	Password             string
	CasePermissions      []CasePermission    `gorm:"foreignKey:UserId"`
	ActionLogs           []ActionLog         `gorm:"foreignKey:UserId"`
	ActorPermissionLogs  []CasePermissionLog `gorm:"foreignKey:ActorUserId"`
	TargetPermissionLogs []CasePermissionLog `gorm:"foreignKey:TargetUserId"`
}
