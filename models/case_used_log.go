package models

import "github.com/google/uuid"

type CaseUsedLog struct {
	Base
	UserId uuid.UUID
	CaseId uuid.UUID
	Count  int `gorm:"default:0"`
}
