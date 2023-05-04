package models

import "github.com/google/uuid"

type File struct {
	Base
	TypeId           *uuid.UUID
	FolderId         *uuid.UUID
	CaseId           *uuid.UUID
	Name             string
	CloudName        string
	PreviewCloudName string
	Url              *FileUrl      `gorm:"-"`
	Tags             []Tag         `gorm:"many2many:file_tags;"`
	ActionLogs       []ActionLog   `gorm:"foreignKey:FileId"`
	ViewLogs         []FileViewLog `gorm:"foreignKey:FileId"`
}
