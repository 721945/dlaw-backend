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
	FileType         *FileType     `gorm:"foreignKey:TypeId"`
	Url              *FileUrl      `gorm:"-"`
	Tags             []Tag         `gorm:"many2many:file_tags;"`
	ActionLogs       []ActionLog   `gorm:"foreignKey:FileId"`
	ViewLogs         []FileViewLog `gorm:"foreignKey:FileId"`
}

type MeiliFile struct {
	Id        string   `json:"id"`
	Name      string   `json:"name,omitempty"`
	Type      string   `json:"type,omitempty"`
	Content   string   `json:"content,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	FolderIds []string `json:"folder_ids,omitempty"`
	CaseId    string   `json:"case_id,omitempty"`
}

type MeiliFileResponse struct {
	ID string `json:"id"`
}
