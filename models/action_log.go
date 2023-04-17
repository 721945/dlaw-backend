package models

import "github.com/google/uuid"

type ActionLog struct {
	Base
	FolderId             uuid.UUID
	FileId               *uuid.UUID
	UserId               uuid.UUID
	ActionId             uuid.UUID
	FilePreviewVersionId *string
	FileVersionId        *string
}
