package models

import "gorm.io/gorm"

type FileType struct {
	gorm.Model
	Name  string
	Files []File `gorm:"foreignKey:TypeId"`
}
