package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type UserRoute struct {
	logger         *libs.Logger
	handler        libs.RequestHandler
	userController controllers.UserController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewUserRoute(
	logger *libs.Logger,
	handler libs.RequestHandler,
	userController controllers.UserController,
	authMiddleware middlewares.JWTAuthMiddleware,
) UserRoute {
	return UserRoute{
		logger:         logger,
		handler:        handler,
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

func (u UserRoute) Setup() {
	u.logger.Info("Setting user routes")
	private := u.handler.Gin.Group("/users")
	{
		private.POST("", u.userController.CreateUser)
		private.GET("", u.userController.GetUsers)
		private.GET("/:id", u.userController.GetUser)
		private.GET("/me", u.authMiddleware.Handler(), u.userController.GetMe)
		private.PATCH("/:id", u.userController.UpdateUser)
		private.PATCH("/me", u.authMiddleware.Handler(), u.userController.UpdateUser)
	}
}
