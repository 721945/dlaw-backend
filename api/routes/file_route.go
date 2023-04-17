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
	api := r.handler.Gin.Group("/files")
	api.Use(r.authMiddleware.Handler())
	{
		api.GET("", r.ctrl.GetFiles)
		api.POST("", r.ctrl.CreateFile)
		api.GET("/:id", r.ctrl.GetFile)
		api.DELETE("", r.ctrl.DeleteFile)
		api.PUT("", r.ctrl.UpdateFile)
		api.POST("/upload", r.ctrl.UploadFile)
	}
}
