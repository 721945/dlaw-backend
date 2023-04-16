package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type AppointmentService struct {
	logger          *libs.Logger
	appointmentRepo repositories.AppointmentRepository
}

func NewAppointmentService(logger *libs.Logger, r repositories.AppointmentRepository) AppointmentService {
	return AppointmentService{logger: logger, appointmentRepo: r}
}

func (s *AppointmentService) GetAppointments() (appointments []models.Appointment, err error) {
	return s.appointmentRepo.GetAppointments()
}

func (s *AppointmentService) GetAppointment(id uuid.UUID) (appointment *models.Appointment, err error) {
	return s.appointmentRepo.GetAppointment(id)
}

func (s *AppointmentService) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepo.CreateAppointment(appointment)
}

func (s *AppointmentService) UpdateAppointment(id uuid.UUID, appointment models.Appointment) error {
	return s.appointmentRepo.UpdateAppointment(id, appointment)
}

func (s *AppointmentService) DeleteAppointment(id uuid.UUID) error {
	return s.appointmentRepo.DeleteAppointment(id)
}
