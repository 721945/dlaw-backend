package models

import (
	"gorm.io/gorm"
)

type Action struct {
	gorm.Model
	Name           string
	ActionLogs     []ActionLog         `gorm:"foreignKey:ActionId"`
	PermissionLogs []CasePermissionLog `gorm:"foreignKey:ActionId"`
}
