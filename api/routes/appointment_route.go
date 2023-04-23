package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type AppointmentRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.AppointmentController
}

func NewAppointmentRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	ctrl controllers.AppointmentController,
) AppointmentRoute {
	return AppointmentRoute{
		handler: handler,
		logger:  logger,
		ctrl:    ctrl,
	}
}

func (r AppointmentRoute) Setup() {
	r.logger.Info("Setting appointment routes")
	api := r.handler.Gin.Group("/appointments")
	{
		api.GET("", r.ctrl.GetAppointments)
		api.POST("", r.ctrl.CreateAppointment)
		api.GET("/:id", r.ctrl.GetAppointment)
		api.DELETE("", r.ctrl.DeleteAppointment)
		api.PUT("", r.ctrl.UpdateAppointment)
		api.GET("/me", r.ctrl.GetOwnAppointment)
		api.GET("/public", r.ctrl.GetOwnAppointment)

	}

}
