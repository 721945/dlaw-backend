package models

type Tag struct {
	Base
	Name     string
	ShowMenu bool `gorm:"default:false"`
}
