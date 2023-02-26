package routes

import "github.com/721945/dlaw-backend/libs"

type FileTypeRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewFileTypeRoute(handler libs.RequestHandler, logger *libs.Logger) FileTypeRoute {
	return FileTypeRoute{handler: handler, logger: logger}
}

func (r FileTypeRoute) Setup() {
	r.logger.Info("Setting file_type routes")
	api := r.handler.Gin.Group("/file_types")
	{
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
	}
}
