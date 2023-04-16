package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
)

type AppointmentRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewAppointmentRepository(logger *libs.Logger, db libs.Database) AppointmentRepository {
	return AppointmentRepository{logger: logger, db: db}
}

func (r *AppointmentRepository) GetAppointments() (appointments []models.Appointment, err error) {
	return appointments, r.db.DB.Find(&appointments).Error

}

func (r *AppointmentRepository) GetAppointment(id uuid.UUID) (appointment *models.Appointment, err error) {
	return appointment, r.db.DB.First(&appointment, id).Error
}

func (r *AppointmentRepository) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return appointment, r.db.DB.Create(&appointment).Error
}

func (r *AppointmentRepository) UpdateAppointment(id uuid.UUID, appointment models.Appointment) error {
	return r.db.DB.Model(&models.Appointment{}).Where("id = ?", id).Updates(appointment).Error
}

func (r *AppointmentRepository) DeleteAppointment(id uuid.UUID) error {
	return r.db.DB.Delete(&models.Appointment{}, id).Error
}
