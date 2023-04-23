package models

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	Base
	CaseId      uuid.UUID
	Case        *Case
	EventId     string
	Title       string
	Detail      string
	Location    string
	DateTime    time.Time
	Emails      []string `gorm:"type:text"`
	IsPublished bool
}
