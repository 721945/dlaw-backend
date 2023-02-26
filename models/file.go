package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Type FileType
	Name string
	Urls []FileUrl `gorm:"many2many:file_urls;"`
	//PublishedUrl string
	Tags []Tag `gorm:"many2many:file_tags;"`
}
