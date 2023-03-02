package models

import "gorm.io/gorm"

type Case struct {
	gorm.Model
	CaseNumber     string
	Title          string
	Folder         Folder `gorm:"foreignKey:CaseId"`
	Email          string
	CaseTitle      string
	CaseDetail     string
	AppointmentId  *uint
	CasePermission []CasePermission    `gorm:"foreignKey:CaseId"`
	PermissionLogs []CasePermissionLog `gorm:"foreignKey:CaseId"`
}
