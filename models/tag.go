package models

import "github.com/google/uuid"

type Tag struct {
	Base
	Name        string
	DisplayName string
	ShowMenu    bool     `gorm:"default:false"`
	Count       int      `gorm:"-"`
	Files       []File   `gorm:"many2many:file_tags;"`
	Folders     []Folder `gorm:"many2many:folder_tags;"`
}

type TagCount struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Count       int
}
