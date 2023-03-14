package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type FileRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.FileController
}

func NewFileRoute(handler libs.RequestHandler, logger *libs.Logger, controller controllers.FileController) FileRoute {
	return FileRoute{handler: handler, logger: logger, ctrl: controller}
}

func (r FileRoute) Setup() {
	r.logger.Info("Setting file routes")
	api := r.handler.Gin.Group("/files")
	{
		api.GET("", r.ctrl.GetFiles)
		api.POST("", r.ctrl.CreateFile)
		api.GET("/:id", r.ctrl.GetFile)
		api.DELETE("", r.ctrl.DeleteFile)
		api.PUT("", r.ctrl.UpdateFile)
		api.POST("/upload", r.ctrl.UploadFile)
	}
}
