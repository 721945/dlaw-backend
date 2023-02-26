package models

import "gorm.io/gorm"

type Action struct {
	gorm.Model
	Name string
}
