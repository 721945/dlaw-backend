package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name           string
	IsArchive      bool
	SubFolders     []Folder `gorm:"foreignKey:ParentFolderId"`
	Files          []File   `gorm:"foreignKey:FolderId"`
	Tags           []Tag    `gorm:"many2many:folder_tags;"`
	ParentFolderId *uint
	CaseId         *uint
	ActionLogs     []ActionLog `gorm:"foreignKey:FolderId"`
}
