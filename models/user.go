package models

import "time"

type User struct {
	Base
	Firstname            string
	Lastname             string
	Email                string
	Password             string
	OtpSecret            *string
	OtpExpiredAt         *time.Time
	Organization         string
	CasePermissions      []CasePermission    `gorm:"foreignKey:UserId"`
	ActionLogs           []ActionLog         `gorm:"foreignKey:UserId"`
	ActorPermissionLogs  []CasePermissionLog `gorm:"foreignKey:ActorUserId"`
	TargetPermissionLogs []CasePermissionLog `gorm:"foreignKey:TargetUserId"`
	CaseUsedLogs         []CaseUsedLog       `gorm:"foreignKey:UserId"`
	FileViewLogs         []FileViewLog       `gorm:"foreignKey:UserId"`
}
