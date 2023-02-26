package routes

import "github.com/721945/dlaw-backend/libs"

type CaseRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewCaseRoute(handler libs.RequestHandler, logger *libs.Logger) CaseRoute {
	return CaseRoute{handler: handler, logger: logger}
}

func (r CaseRoute) Setup() {
	r.logger.Info("Setting case routes")
	api := r.handler.Gin.Group("/cases")
	{
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
	}
}
