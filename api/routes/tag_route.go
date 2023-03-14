package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type TagRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.TagController
}

func NewTagRoute(handler libs.RequestHandler, logger *libs.Logger, ctrl controllers.TagController) TagRoute {
	return TagRoute{handler: handler, logger: logger, ctrl: ctrl}
}

func (r TagRoute) Setup() {
	r.logger.Info("Setting tag routes")
	api := r.handler.Gin.Group("/tags")
	{
		api.GET("", r.ctrl.GetTags)
		api.POST("", r.ctrl.CreateTag)
		api.GET("/:id", r.ctrl.GetTag)
		api.DELETE("/:id", r.ctrl.DeleteTag)
		api.PUT("/:id", r.ctrl.UpdateTag)
	}
}
