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
		api.GET("/me", r.caseCtrl.GetOwnCases)
		api.GET("/freq", r.caseCtrl.GetOwnCases)

		//api.DELETE("", r.caseCtrl.DeleteCase)
		//api.PUT("", r.caseCtrl.UpdateCase)
	}
}
