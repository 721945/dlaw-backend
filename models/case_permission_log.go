package models

import "gorm.io/gorm"

type CasePermissionLog struct {
	gorm.Model
	PermissionId uint
	ActionId     uint
	ActorUserId  uint
	TargetUserId uint
	CaseId       uint
}
