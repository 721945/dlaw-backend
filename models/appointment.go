package models

type Appointment struct {
	Base
	Case        Case
	EventId     string
	Title       string
	Detail      string
	IsPublished bool
}
