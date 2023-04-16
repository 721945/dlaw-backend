package models

type FileType struct {
	Base
	Name  string
	Files []File `gorm:"foreignKey:TypeId"`
}
