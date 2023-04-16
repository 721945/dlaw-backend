package models

type User struct {
	Base
	Firstname            string
	Lastname             string
	Email                string
	Password             string
	CasePermissions      []CasePermission    `gorm:"foreignKey:UserId"`
	ActionLogs           []ActionLog         `gorm:"foreignKey:UserId"`
	ActorPermissionLogs  []CasePermissionLog `gorm:"foreignKey:ActorUserId"`
	TargetPermissionLogs []CasePermissionLog `gorm:"foreignKey:TargetUserId"`
}
