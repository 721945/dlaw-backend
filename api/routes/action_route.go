package routes

import "github.com/721945/dlaw-backend/libs"

type ActionRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewActionRoute(handler libs.RequestHandler, logger *libs.Logger) ActionRoute {
	return ActionRoute{handler: handler, logger: logger}
}

func (r ActionRoute) Setup() {
	r.logger.Info("Setting action routes")
	api := r.handler.Gin.Group("/actions")
	{
		api.GET("", nil)
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)

	}

}
