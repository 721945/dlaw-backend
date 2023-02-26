package routes

import "github.com/721945/dlaw-backend/libs"

type FolderRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewFolderRoute(handler libs.RequestHandler, logger *libs.Logger) FolderRoute {
	return FolderRoute{handler: handler, logger: logger}
}

func (r FolderRoute) Setup() {
	r.logger.Info("Setting folder routes")
	api := r.handler.Gin.Group("/folders")
	{
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
	}
}
