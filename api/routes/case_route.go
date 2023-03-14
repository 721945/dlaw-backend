package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type CaseRoute struct {
	handler  libs.RequestHandler
	logger   *libs.Logger
	caseCtrl controllers.CaseController
}

func NewCaseRoute(handler libs.RequestHandler, logger *libs.Logger, ctrl controllers.CaseController) CaseRoute {
	return CaseRoute{handler: handler, logger: logger, caseCtrl: ctrl}
}

func (r CaseRoute) Setup() {
	r.logger.Info("Setting case routes")
	api := r.handler.Gin.Group("/cases")
	{
		api.GET("", r.caseCtrl.GetCases)
		api.POST("", r.caseCtrl.CreateCase)
		api.GET("/:id", r.caseCtrl.GetCase)

		//api.DELETE("", r.caseCtrl.DeleteCase)
		//api.PUT("", r.caseCtrl.UpdateCase)
	}
}
