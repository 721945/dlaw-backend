package routes

import "github.com/721945/dlaw-backend/libs"

type AppointmentRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewAppointmentRoute(handler libs.RequestHandler, logger *libs.Logger) AppointmentRoute {
	return AppointmentRoute{handler: handler, logger: logger}
}

func (r AppointmentRoute) Setup() {
	r.logger.Info("Setting appointment routes")
	api := r.handler.Gin.Group("/appointments")
	{
		api.GET("", nil)
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)

	}

}
