package models

import "github.com/google/uuid"

type CaseUsedLog struct {
	Base
	UserId uuid.UUID
	CaseId uuid.UUID
	count  int
}
