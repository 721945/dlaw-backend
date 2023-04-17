package models

type Case struct {
	Base
	RedCaseNumber   string
	BlackCaseNumber string
	Title           string
	CaseTitle       *string
	CaseDetail      *string
	Email           *string
	AppointmentId   *uint
	Folders         []Folder            `gorm:"foreignKey:CaseId"`
	CasePermission  []CasePermission    `gorm:"foreignKey:CaseId"`
	PermissionLogs  []CasePermissionLog `gorm:"foreignKey:CaseId"`
}
