package models

import "gorm.io/gorm"

type Appointment struct {
	gorm.Model
	Case        Case
	EventId     string
	Title       string
	Detail      string
	IsPublished bool
}
