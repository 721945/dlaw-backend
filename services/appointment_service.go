package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/infrastructure/google_calendar"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type AppointmentService struct {
	logger          *libs.Logger
	appointmentRepo repositories.AppointmentRepository
	caseRepo        repositories.CaseRepository
	casePermission  repositories.CasePermissionRepository
	calendarService google_calendar.GoogleCalendar
}

func NewAppointmentService(
	logger *libs.Logger,
	r repositories.AppointmentRepository,
	cr repositories.CaseRepository,
	cpr repositories.CasePermissionRepository,
	cld google_calendar.GoogleCalendar,
) AppointmentService {
	return AppointmentService{
		logger:          logger,
		appointmentRepo: r,
		caseRepo:        cr,
		casePermission:  cpr,
		calendarService: cld,
	}
}

func (s *AppointmentService) GetAppointments() (appointments []dtos.AppointmentDto, err error) {
	appointmentModels, err := s.appointmentRepo.GetAppointments()

	if err != nil {
		return nil, err
	}

	appointments = dtos.ToAppointmentDtoList(appointmentModels)

	return appointments, nil
}

func (s *AppointmentService) GetAppointmentsByUser(userId uuid.UUID) (appointments []dtos.AppointmentDto, err error) {

	caseIds, err := s.getCaseIds(userId)

	appointmentModels, err := s.appointmentRepo.GetAppointmentByCaseIds(caseIds)

	if err != nil {
		return nil, err
	}

	appointments = dtos.ToAppointmentDtoList(appointmentModels)

	return
}

func (s *AppointmentService) GetAppointment(id uuid.UUID) (appointment *dtos.AppointmentDto, err error) {
	model, err := s.appointmentRepo.GetAppointment(id)

	if err != nil {
		return nil, err
	}

	dto := dtos.ToAppointmentDto(*model)

	return &dto, nil
}

func (s *AppointmentService) CreateAppointment(userId uuid.UUID, dto dtos.CreateAppointmentDto) (string, error) {
	id, err := uuid.Parse(dto.CaseId)
	if err != nil {
		return "", err
	}

	// Create google calendar event
	caseId, err := s.getCaseByUserIdAndCaseId(userId, id)

	if err != nil {
		return "", err
	}

	if caseId == nil {
		return "", libs.ErrNotFound
	}

	// Create event
	// event, err := s.calendarService.CreateEvent(dto.Title, dto.DateTime, dto.DateTime, dto.Emails, &dto.Location, &dto.Detail)
	if err != nil {
		return "", err
	}

	appointmentModel := dto.ToAppointmentModel(*caseId, "")

	appointment, err := s.appointmentRepo.CreateAppointment(appointmentModel)

	return appointment.ID.String(), err
}

func (s *AppointmentService) UpdateAppointment(id uuid.UUID, dto dtos.UpdateAppointmentDto, userId uuid.UUID) error {
	caseId, err := s.getCaseByUserIdAndCaseId(userId, id)

	if err != nil {
		return err
	}

	if caseId == nil {
		return libs.ErrNotFound
	}

	appointmentModel := dto.ToAppointmentModel()

	return s.appointmentRepo.UpdateAppointment(id, appointmentModel)
}

func (s *AppointmentService) DeleteAppointment(id uuid.UUID, userId uuid.UUID) error {
	caseId, err := s.getCaseByUserIdAndCaseId(userId, id)

	if err != nil {
		return err
	}

	if caseId == nil {
		return libs.ErrNotFound
	}

	return s.appointmentRepo.DeleteAppointment(id)
}

func (s *AppointmentService) getCaseIds(userId uuid.UUID) ([]uuid.UUID, error) {
	cases, err := s.casePermission.GetCasePermissionsByUserId(userId)
	if err != nil {
		return nil, err
	}

	caseIds := make([]uuid.UUID, len(cases))
	for i, c := range cases {
		caseIds[i] = c.CaseId
	}

	return caseIds, nil
}

func (s *AppointmentService) getCaseByUserIdAndCaseId(userId uuid.UUID, caseId uuid.UUID) (*uuid.UUID, error) {
	caseModel, err := s.casePermission.GetCasePermissionsByUserIdAndCaseId(userId, caseId)
	if err != nil {
		return nil, err
	}

	return &caseModel.CaseId, nil
}

func (s *AppointmentService) GetPublicAppointment() (appointments []dtos.AppointmentDto, err error) {
	appointmentModels, err := s.appointmentRepo.GetPublicAppointments()

	if err != nil {
		return nil, err
	}

	appointments = dtos.ToAppointmentDtoList(appointmentModels)

	return appointments, err
}
