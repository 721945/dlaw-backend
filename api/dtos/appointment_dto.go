package dtos

import (
	"fmt"
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type CreateAppointmentDto struct {
	CaseId   string   `json:"caseId" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Location string   `json:"location" `
	Detail   string   `json:"detail" `
	Emails   []string `json:"emails" binding:"required"`
	//StartDate string   `json:"startDate" binding:"required"`
	DateTime string `json:"dateTime" binding:"required"`
}

type UpdateAppointmentDto struct {
	Location    string   `json:"location" `
	Detail      string   `json:"detail" `
	Title       string   `json:"title" `
	Emails      []string `json:"emails" `
	DateTime    string   `json:"dateTime" `
	IsPublished bool     `json:"isPublished" `
}

type AppointmentDto struct {
	ID          string   `json:"id"`
	Emails      []string `json:"emails"`
	Title       string   `json:"title"`
	Detail      string   `json:"detail"`
	Location    string   `json:"location"`
	DateTime    string   `json:"dateTime"`
	UpdatedAt   string   `json:"updatedAt"`
	IsPublished bool     `json:"isPublished"`
}

func (c *CreateAppointmentDto) ToAppointmentModel(caseId uuid.UUID, eventId string) models.Appointment {
	t, err := utils.CovertStringToTime(c.DateTime)
	if err != nil {
		fmt.Println(err)
		return models.Appointment{}
	}
	emails := make([]models.Email, len(c.Emails))
	for i, email := range c.Emails {
		emails[i] = models.Email{
			Email: email,
		}
	}
	return models.Appointment{
		CaseId:      caseId,
		EventId:     eventId,
		Title:       c.Title,
		Detail:      c.Detail,
		Location:    c.Location,
		DateTime:    t,
		Emails:      emails,
		IsPublished: false,
	}
}

func (u *UpdateAppointmentDto) ToAppointmentModel() models.Appointment {

	t, err := utils.CovertStringToTime(u.DateTime)
	if err != nil {
		fmt.Println(err)
		return models.Appointment{}
	}
	emails := make([]models.Email, len(u.Emails))
	for i, email := range u.Emails {
		emails[i] = models.Email{
			Email: email,
		}
	}

	return models.Appointment{
		Title:    u.Title,
		Location: u.Location,
		Detail:   u.Detail,
		Emails:   emails,
		DateTime: t,
		//DateTime: u.DateTime,
	}
}

func ToAppointmentDto(appointment models.Appointment) AppointmentDto {
	updatedAt := utils.CovertTimeToString(appointment.UpdatedAt)
	if appointment.CreatedAt == appointment.UpdatedAt {
		updatedAt = ""
	}
	emails := make([]string, len(appointment.Emails))
	for i, email := range appointment.Emails {
		emails[i] = email.Email
	}

	return AppointmentDto{
		ID:          appointment.ID.String(),
		Emails:      emails,
		Title:       appointment.Title,
		Detail:      appointment.Detail,
		Location:    appointment.Location,
		DateTime:    utils.CovertTimeToString(appointment.DateTime),
		UpdatedAt:   updatedAt,
		IsPublished: appointment.IsPublished,
	}
}

func ToAppointmentDtoList(appointments []models.Appointment) []AppointmentDto {
	appointmentDtos := make([]AppointmentDto, len(appointments))
	for i, appointment := range appointments {
		appointmentDtos[i] = ToAppointmentDto(appointment)
	}
	return appointmentDtos
}
