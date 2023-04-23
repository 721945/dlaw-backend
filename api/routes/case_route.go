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
		api.GET("/:id", r.caseCtrl.GetCase)
		api.GET("/me", r.caseCtrl.GetOwnCases)
		api.GET("/freq", r.caseCtrl.GetOwnCases)
		api.GET("/archived", r.caseCtrl.GetArchivedCases)
		api.GET("/:id/members", r.caseCtrl.GetMembers)
		api.POST("", r.caseCtrl.CreateCase)
		api.POST("/:id/members", r.caseCtrl.AddMember)
		api.PATCH("/:id/members/:member", r.caseCtrl.UpdateMember)
		api.PATCH("/:id", r.caseCtrl.UpdateCase)
		api.PATCH("/:id/archive", r.caseCtrl.ArchiveCase)
		api.DELETE("/:id", r.caseCtrl.DeleteCase)
		api.DELETE("/:id/members/:member", r.caseCtrl.RemoveMember)
		//api.PUT("", r.caseCtrl.UpdateCase)
	}
}
