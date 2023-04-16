package models

import "github.com/google/uuid"

type Folder struct {
	Base
	Name           string
	IsArchive      bool
	ParentFolderId *uuid.UUID
	CaseId         *uuid.UUID
	SubFolders     []Folder    `gorm:"foreignKey:ParentFolderId"`
	Files          []File      `gorm:"foreignKey:FolderId"`
	Tags           []Tag       `gorm:"many2many:folder_tags;"`
	ActionLogs     []ActionLog `gorm:"foreignKey:FolderId"`
}
