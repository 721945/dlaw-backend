package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type PermissionRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	ctrl    controllers.PermissionController
}

func NewPermissionRoute(handler libs.RequestHandler, logger *libs.Logger, ctrl controllers.PermissionController) PermissionRoute {
	return PermissionRoute{handler: handler, logger: logger, ctrl: ctrl}
}

func (r PermissionRoute) Setup() {
	r.logger.Info("Setting permission routes")
	api := r.handler.Gin.Group("/permissions")
	{
		api.GET("", r.ctrl.GetPermissions)
		api.GET("/:id", r.ctrl.GetPermission)
		api.GET("/name/:name", r.ctrl.GetPermissionName)
		api.POST("", r.ctrl.CreatePermission)
		api.DELETE("/:id", r.ctrl.DeletePermission)
		api.PUT("/:id", r.ctrl.UpdatePermission)
	}
}
