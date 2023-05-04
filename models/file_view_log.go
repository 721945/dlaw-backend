package models

import "github.com/google/uuid"

type FileViewLog struct {
	Base
	FileId uuid.UUID
	File   *File
	UserId uuid.UUID
}
