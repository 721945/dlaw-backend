package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type FolderRoute struct {
	handler    libs.RequestHandler
	logger     *libs.Logger
	controller controllers.FolderController
}

func NewFolderRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	controller controllers.FolderController,
) FolderRoute {
	return FolderRoute{
		handler:    handler,
		logger:     logger,
		controller: controller,
	}
}

func (r FolderRoute) Setup() {
	r.logger.Info("Setting folder routes")
	api := r.handler.Gin.Group("/folders")
	{
		api.GET("", r.controller.GetFolders)
		api.GET("/:id", r.controller.GetFolder)
		api.GET("/:id/root", r.controller.GetFolderRoot)
		api.GET("/:id/log", r.controller.GetFolderLog)

		api.POST("", r.controller.CreateFolder)

		api.PATCH("/:id/move", r.controller.MoveFolder)
		api.PATCH("/:id", r.controller.UpdateFolder)

		api.DELETE("/:id", r.controller.DeleteFolder)
	}
}
