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
	public := r.handler.Gin.Group("/appointments")
	{
		public.GET("/public", r.ctrl.GetPublicAppointment)
	}
	
	private := r.handler.Gin.Group("/appointments")
	private.Use(r.authMiddleware.Handler())
	{
		private.GET("", r.ctrl.GetAppointments)
		private.POST("", r.ctrl.CreateAppointment)
		private.GET("/:id", r.ctrl.GetAppointment)
		private.DELETE("/:id", r.ctrl.DeleteAppointment)
		private.PATCH("/:id", r.ctrl.UpdateAppointment)
		private.PATCH("/:id/publish", r.ctrl.PublishedAppointment)
		private.PATCH("/:id/un-publish", r.ctrl.UnPublishedAppointment)
		private.GET("/me", r.ctrl.GetOwnAppointment)
	}

}
