package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type UserRoutes struct {
	logger         libs.Logger
	handler        libs.RequestHandler
	userController controllers.UserController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewUserRoutes(
	logger libs.Logger,
	handler libs.RequestHandler,
	userController controllers.UserController,
	authMiddleware middlewares.JWTAuthMiddleware,
) UserRoutes {
	return UserRoutes{logger: logger, handler: handler, userController: userController, authMiddleware: authMiddleware}
}

func (u UserRoutes) Setup() {
	u.logger.Info("Setting user routes")
	api := u.handler.Gin.Group("/users")
	//.Use(u.authMiddleware.Handler())
	{
		api.POST("", u.userController.CreateUser)
		api.GET("", u.userController.GetUsers)
		api.GET("/:id", u.userController.GetUser)
		api.GET("/me", u.userController.GetMe).Use(u.authMiddleware.Handler())

		//api.GET("/ping", u.userController.Ping)
	}
}
