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

	public := r.handler.Gin.Group("/cases")
	{
		public.GET("/public", r.caseCtrl.GetPublicCases)
	}

	private := r.handler.Gin.Group("/cases")
	private.Use(r.authMiddleware.Handler())
	{
		private.GET("", r.caseCtrl.GetCases)
		private.GET("/:id", r.caseCtrl.GetCase)
		private.GET("/me", r.caseCtrl.GetOwnCases)
		private.GET("/freq", r.caseCtrl.GetFrequentlyUsed)
		private.GET("/archived", r.caseCtrl.GetArchivedCases)
		private.GET("/:id/members", r.caseCtrl.GetMembers)
		private.GET("/:id/folders", r.caseCtrl.GetFolders)
		private.POST("", r.caseCtrl.CreateCase)
		private.POST("/:id/members", r.caseCtrl.AddMember)
		private.PATCH("/:id/members/:member", r.caseCtrl.UpdateMember)
		private.PATCH("/:id", r.caseCtrl.UpdateCase)
		private.PATCH("/:id/archive", r.caseCtrl.ArchiveCase)
		private.DELETE("/:id", r.caseCtrl.DeleteCase)
		private.DELETE("/:id/members/:member", r.caseCtrl.RemoveMember)
	}
}

//func DebugMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var names []string
//		for _, h := range c.HandlerNames() {
//			names = append(names, h)
//		}
//		fmt.Println("Middleware functions:", strings.Join(names, ", "))
//		c.Next()
//	}
//}
