package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type FileTypeRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.TypeController
}

func NewFileTypeRoute(handler libs.RequestHandler, logger *libs.Logger, typeCtrl controllers.TypeController) FileTypeRoute {
	return FileTypeRoute{handler: handler, logger: logger, ctrl: typeCtrl}
}

func (r FileTypeRoute) Setup() {
	r.logger.Info("Setting file_type routes")
	api := r.handler.Gin.Group("/file_types")
	{
		api.GET("", r.ctrl.GetTypes)
		api.POST("", r.ctrl.CreateType)
		api.POST("/all", r.ctrl.CreateTypes)
		api.GET("/:id", r.ctrl.GetType)
		api.DELETE("/:id", r.ctrl.DeleteType)
		api.PUT("/:id", r.ctrl.UpdateType)
	}
}
