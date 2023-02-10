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

func NewUserRoutes(logger libs.Logger,
	handler libs.RequestHandler,
	userController controllers.UserController,
) UserRoutes {
	return UserRoutes{
		logger:         logger,
		handler:        handler,
		userController: userController,
	}
}

func (u UserRoutes) Setup() {
	u.logger.Info("Setting user routes")
	//api := u.handler.Gin.Group("/user").Use(u.authMiddleware.Handler())
	{
		//api.GET("/ping", u.userController.Ping)
	}
}
