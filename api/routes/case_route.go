package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type CaseRoute struct {
	handler        libs.RequestHandler
	logger         *libs.Logger
	caseCtrl       controllers.CaseController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewCaseRoute(
	handler libs.RequestHandler,
	logger *libs.Logger,
	caseCtrl controllers.CaseController,
	authMiddleware middlewares.JWTAuthMiddleware,
) CaseRoute {
	return CaseRoute{
		handler:        handler,
		logger:         logger,
		caseCtrl:       caseCtrl,
		authMiddleware: authMiddleware,
	}
}

func (r CaseRoute) Setup() {
	r.logger.Info("Setting case routes")
	api := r.handler.Gin.Group("/cases")
	api.Use(r.authMiddleware.Handler())
	{
		api.GET("", r.caseCtrl.GetCases)
		api.POST("", r.caseCtrl.CreateCase)
		api.GET("/:id", r.caseCtrl.GetCase)
		api.GET("/archived", r.caseCtrl.GetArchivedCases)
		api.GET("/me", r.caseCtrl.GetOwnCases)
		api.GET("/freq", r.caseCtrl.GetOwnCases)
		api.PATCH("/:id", r.caseCtrl.UpdateCase)
		api.PATCH("/:id/archive", r.caseCtrl.ArchiveCase)
		//api.PUT("", r.caseCtrl.UpdateCase)
	}
}
