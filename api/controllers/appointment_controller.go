package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	logger             *libs.Logger
	appointmentService services.AppointmentService
}

func NewAppointmentController(logger *libs.Logger, appointmentService services.AppointmentService) AppointmentController {
	return AppointmentController{logger: logger, appointmentService: appointmentService}
}

func (p AppointmentController) GetAppointments(c *gin.Context) {
	permissions, err := p.appointmentService.GetAppointments()
	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permissions})
}

func (p AppointmentController) GetOwnAppointment(c *gin.Context) {

	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	permissions, err := p.appointmentService.GetAppointmentsByUser(userId.(uuid.UUID))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permissions})
}

func (p AppointmentController) GetAppointment(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)
	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	permission, err := p.appointmentService.GetAppointment(id)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permission})
}

//func (p AppointmentController) GetAppointmentName(c *gin.Context) {
//	name := c.Param("name")
//
//	permission, err := p.appointmentService.GetAppointmentByName(name)
//
//	if err != nil {
//		p.logger.Error(err)
//		_ = c.Error(err)
//		return
//	}
//
//	c.JSON(200, gin.H{"data": permission})
//}

func (p AppointmentController) CreateAppointment(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	var input dtos.CreateAppointmentDto

	if err := c.ShouldBindJSON(&input); err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	permission, err := p.appointmentService.CreateAppointment(userId.(uuid.UUID), input)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(201, gin.H{"data": permission})
}

func (p AppointmentController) UpdateAppointment(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.UpdateAppointmentDto
	if err := c.ShouldBindJSON(&input); err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = p.appointmentService.UpdateAppointment(id, input, userId.(uuid.UUID))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})

}

func (p AppointmentController) PublishedAppointment(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.UpdateAppointmentDto

	input.IsPublished = true

	err = p.appointmentService.UpdateAppointment(id, input, userId.(uuid.UUID))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})

}

func (p AppointmentController) UnPublishedAppointment(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.UpdateAppointmentDto

	input.IsPublished = false

	err = p.appointmentService.UpdateAppointment(id, input, userId.(uuid.UUID))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})

}

func (p AppointmentController) DeleteAppointment(c *gin.Context) {

	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = p.appointmentService.DeleteAppointment(id, userId.(uuid.UUID))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

func (p AppointmentController) GetPublicAppointment(c *gin.Context) {

	permissions, err := p.appointmentService.GetPublicAppointment()

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permissions})
}
