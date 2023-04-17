package models

import "github.com/google/uuid"

type CasePermissionLog struct {
	Base
	PermissionId uuid.UUID
	ActionId     uuid.UUID
	ActorUserId  uuid.UUID
	TargetUserId uuid.UUID
	CaseId       uuid.UUID
}
