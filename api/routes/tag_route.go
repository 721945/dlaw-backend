package routes

import "github.com/721945/dlaw-backend/libs"

type TagRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewTagRoute(handler libs.RequestHandler, logger *libs.Logger) TagRoute {
	return TagRoute{handler: handler, logger: logger}
}

func (r TagRoute) Setup() {
	r.logger.Info("Setting tag routes")
	api := r.handler.Gin.Group("/tags")
	{
		api.GET("", nil)
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
	}
}
