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

type MeiliFile struct {
	Id        string   `json:"id"`
	FileId    string   `json:"file_id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	FolderIds []string `json:"folder_ids"`
	CaseId    string   `json:"case_id"`
}
