package models

import "gorm.io/gorm"

type FileUrl struct {
	gorm.Model
	Url          string
	PublishedUrl string
	FileId       *uint
}
