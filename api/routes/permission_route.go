package routes

import "github.com/721945/dlaw-backend/libs"

type PermissionRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewPermissionRoute(handler libs.RequestHandler, logger *libs.Logger) PermissionRoute {
	return PermissionRoute{handler: handler, logger: logger}
}

func (r PermissionRoute) Setup() {
	r.logger.Info("Setting permission routes")
	api := r.handler.Gin.Group("/permissions")
	{
		api.POST("", nil)
		api.GET("/:id", nil)
		api.DELETE("", nil)
		api.PUT("", nil)
	}
}
