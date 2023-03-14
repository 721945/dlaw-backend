package models

import "gorm.io/gorm"

type ActionLog struct {
	gorm.Model
	FolderId uint
	FileId   *uint
	UserId   uint
	ActionId uint
}
