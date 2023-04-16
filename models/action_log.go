package models

import "github.com/google/uuid"

type ActionLog struct {
	Base
	FolderId uuid.UUID
	FileId   *uint
	UserId   uint
	ActionId uuid.UUID
}
