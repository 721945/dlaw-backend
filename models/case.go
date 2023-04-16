package models

type Case struct {
	Base
	RedCaseNumber   string
	BlackCaseNumber string
	Title           string
	Folder          Folder `gorm:"foreignKey:CaseId"`
	CaseTitle       *string
	CaseDetail      *string
	Email           *string
	AppointmentId   *uint
	CasePermission  []CasePermission    `gorm:"foreignKey:CaseId"`
	PermissionLogs  []CasePermissionLog `gorm:"foreignKey:CaseId"`
}
