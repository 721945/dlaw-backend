package models

import "github.com/google/uuid"

type ActionLog struct {
	Base
	FolderId             uuid.UUID
	Folder               *Folder
	FileId               *uuid.UUID
	File                 *File
	UserId               uuid.UUID
	User                 *User
	ActionId             uuid.UUID
	Action               *Action
	FilePreviewVersionId *string
	FileVersionId        *string
}
