package routes

import "github.com/721945/dlaw-backend/libs"

type FileRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewFileRoute(handler libs.RequestHandler, logger *libs.Logger) FileRoute {
	return FileRoute{handler: handler, logger: logger}
}

func (r FileRoute) Setup() {
	r.logger.Info("Setting file routes")
	api := r.handler.Gin.Group("/files")
	{
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
		api.POST("/upload", nil)
	}
}
