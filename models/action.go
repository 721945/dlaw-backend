package models

type Action struct {
	Base
	Name           string
	ActionLogs     []ActionLog         `gorm:"foreignKey:ActionId"`
	PermissionLogs []CasePermissionLog `gorm:"foreignKey:ActionId"`
}
