package models

import "github.com/google/uuid"

type File struct {
	Base
	TypeId *uuid.UUID
	Name   string
	Urls   []FileUrl `gorm:"foreignKey:FileId"`
	//PublishedUrl string
	Tags       []Tag `gorm:"many2many:file_tags;"`
	FolderId   *uuid.UUID
	ActionLogs []ActionLog `gorm:"foreignKey:FileId"`
	//UploadBy   User        `gorm:"foreignKey:UploadByUserId"`
}
