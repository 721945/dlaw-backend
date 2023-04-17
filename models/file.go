package models

import "github.com/google/uuid"

type File struct {
	Base
	TypeId           *uuid.UUID
	FolderId         *uuid.UUID
	Name             string
	CloudName        string
	PreviewCloudName string
	Tags             []Tag       `gorm:"many2many:file_tags;"`
	ActionLogs       []ActionLog `gorm:"foreignKey:FileId"`
	Url              *FileUrl
	//Urls         []FileUrl     `gorm:"foreignKey:FileId"`
	//FileVersions []FileVersion `gorm:"foreignKey:FileId"`
	//UploadBy   User        `gorm:"foreignKey:UploadByUserId"`
}
