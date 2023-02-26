package models

import "gorm.io/gorm"

type ActionLog struct {
	gorm.Model
	Folder Folder
	File   File
	User   User
	Action Action
}
