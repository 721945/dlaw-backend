package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type AppointmentRoute struct {
	handler        libs.RequestHandler
	logger         *libs.Logger
	ctrl           controllers.AppointmentController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewAppointmentRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	ctrl controllers.AppointmentController,
	authMiddleware middlewares.JWTAuthMiddleware,
) AppointmentRoute {
	return AppointmentRoute{
		handler:        handler,
		logger:         logger,
		ctrl:           ctrl,
		authMiddleware: authMiddleware,
	}
}

func (r AppointmentRoute) Setup() {
	r.logger.Info("Setting appointment routes")
	api := r.handler.Gin.Group("/appointments")
	{
		api.GET("", r.ctrl.GetAppointments)
		api.POST("", r.ctrl.CreateAppointment).Use(r.authMiddleware.Handler())
		api.GET("/:id", r.ctrl.GetAppointment).Use(r.authMiddleware.Handler())
		api.DELETE("/:id", r.ctrl.DeleteAppointment).Use(r.authMiddleware.Handler())
		api.PATCH("/:id", r.ctrl.UpdateAppointment).Use(r.authMiddleware.Handler())
		api.GET("/me", r.ctrl.GetOwnAppointment).Use(r.authMiddleware.Handler())
		api.GET("/public", r.ctrl.GetOwnAppointment)

	}

}
