package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type ActionRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.ActionController
}

func NewActionRoute(handler libs.RequestHandler, logger *libs.Logger, ctrl controllers.ActionController) ActionRoute {
	return ActionRoute{handler: handler, logger: logger, ctrl: ctrl}
}

func (r ActionRoute) Setup() {
	r.logger.Info("Setting action routes")
	api := r.handler.Gin.Group("/actions")
	{
		api.GET("", r.ctrl.GetActions)
		api.POST("", r.ctrl.CreateAction)
		api.POST("/all", r.ctrl.CreateActions)
		api.GET("/:id", r.ctrl.GetAction)
		api.DELETE("", r.ctrl.DeleteAction)
		api.PUT("", r.ctrl.UpdateAction)
	}

}
