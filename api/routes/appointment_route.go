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
	api.Use(r.authMiddleware.Handler())
	{

		api.GET("", r.ctrl.GetAppointments)
		api.POST("", r.ctrl.CreateAppointment)
		api.GET("/:id", r.ctrl.GetAppointment)
		api.DELETE("/:id", r.ctrl.DeleteAppointment)
		api.PATCH("/:id", r.ctrl.UpdateAppointment)
		api.PATCH("/:id/publish", r.ctrl.PublishedAppointment)
		api.PATCH("/:id/un-publish", r.ctrl.UnPublishedAppointment)
		api.GET("/me", r.ctrl.GetOwnAppointment)
		api.GET("/public", r.ctrl.GetPublicAppointment)

	}

}
