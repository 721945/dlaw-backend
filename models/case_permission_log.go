package models

import "gorm.io/gorm"

type FolderPermissionLog struct {
	gorm.Model
	Permission Permission
	Action     Action
	User       User
	Case       Case
}
