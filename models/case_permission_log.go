package models

import "github.com/google/uuid"

type CasePermissionLog struct {
	Base
	PermissionId uuid.UUID
	ActionId     uint
	ActorUserId  uint
	TargetUserId uuid.UUID
	CaseId       uint
}
