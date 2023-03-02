package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	TypeId       *uint
	Name         string
	Urls         []FileUrl `gorm:"foreignKey:FileId"`
	PublishedUrl string
	Tags         []Tag `gorm:"many2many:file_tags;"`
	FolderId     *uint
	ActionLogs   []ActionLog `gorm:"foreignKey:FileId"`
}
