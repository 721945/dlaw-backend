package models

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	Base
	CaseId      uuid.UUID
	Case        *Case `gorm:"foreignKey:CaseId"`
	EventId     string
	Title       string
	Detail      string
	Location    string
	DateTime    time.Time
	Emails      []Email
	IsPublished bool `gorm:"default:false"`
}

type Email struct {
	Base
	AppointmentId uuid.UUID
	Appointment   *Appointment
	Email         string
}

//func (a Appointment) BeforeCreate(tx *gorm.DB) (err error) {
//
//	a.Emails = strings.Join(a.Emails, ",")
//
//	return
//}
