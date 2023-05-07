package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type FolderRoute struct {
	handler       libs.RequestHandler
	logger        *libs.Logger
	controller    controllers.FolderController
	jwtMiddleWare middlewares.JWTAuthMiddleware
}

func NewFolderRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	controller controllers.FolderController,
	jwtMiddleWare middlewares.JWTAuthMiddleware,
) FolderRoute {
	return FolderRoute{
		handler:       handler,
		logger:        logger,
		controller:    controller,
		jwtMiddleWare: jwtMiddleWare,
	}
}

func (r FolderRoute) Setup() {
	r.logger.Info("Setting folder routes")
	private := r.handler.Gin.Group("/folders")
	private.Use(r.jwtMiddleWare.Handler())
	{
		private.GET("", r.controller.GetFolders)
		private.GET("/:id", r.controller.GetFolder)
		private.GET("/:id/root", r.controller.GetFolderRoot)
		private.GET("/:id/log", r.controller.GetFolderLog)
		private.GET("/:id/menu", r.controller.GetTagMenus)
		private.GET("/:id/tag/:tagId/files", r.controller.GetFilesInTag)

		private.POST("", r.controller.CreateFolder)

		private.PATCH("/:id/move", r.controller.MoveFolder)
		private.PATCH("/:id", r.controller.UpdateFolder)

		private.DELETE("/:id", r.controller.DeleteFolder)

	}
}
