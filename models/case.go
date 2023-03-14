package models

import "gorm.io/gorm"

type Case struct {
	gorm.Model
	CaseNumber     string
	Title          string
	Folder         Folder `gorm:"foreignKey:CaseId"`
	CaseTitle      *string
	CaseDetail     *string
	Email          *string
	AppointmentId  *uint
	CasePermission []CasePermission    `gorm:"foreignKey:CaseId"`
	PermissionLogs []CasePermissionLog `gorm:"foreignKey:CaseId"`
}
