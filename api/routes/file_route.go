package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type FileRoute struct {
	handler        libs.RequestHandler
	logger         *libs.Logger
	ctrl           controllers.FileController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewFileRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	controller controllers.FileController,
	authMiddleware middlewares.JWTAuthMiddleware,
) FileRoute {
	return FileRoute{
		handler:        handler,
		logger:         logger,
		ctrl:           controller,
		authMiddleware: authMiddleware,
	}
}

func (r FileRoute) Setup() {
	r.logger.Info("Setting file routes")
	public := r.handler.Gin.Group("/files")
	{
		public.GET("/public/:id", r.ctrl.GetPublicFile)
	}
	private := r.handler.Gin.Group("/files")
	private.Use(r.authMiddleware.Handler())
	{
		private.GET("", r.ctrl.GetFiles)
		private.GET("/:id", r.ctrl.GetFile)
		private.GET("/tags/count", r.ctrl.CountFileInTags)
		private.GET("/search/:word", r.ctrl.SearchFiles)
		private.GET("/recent", r.ctrl.RecentViewedFiles)
		private.POST("", r.ctrl.CreateFile)
		private.POST("/upload", r.ctrl.UploadFile)
		private.PATCH("/:id/move", r.ctrl.MoveFile)
		private.PATCH("/:id", r.ctrl.UpdateFile)
		private.PATCH("/:id/share", r.ctrl.ShareFile)
		private.PATCH("/:id/remove-share", r.ctrl.UnShareFile)
		private.PATCH("/:id/public", r.ctrl.PublicFile)
		private.DELETE("/:id", r.ctrl.DeleteFile)
	}
}
