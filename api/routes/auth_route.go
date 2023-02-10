package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/libs"
)

type AuthRoutes struct {
	logger         libs.Logger
	handler        libs.RequestHandler
	authController controllers.AuthController
}

func NewAuthRoutes(
	logger libs.Logger,
	handler libs.RequestHandler,
	authController controllers.AuthController,
) AuthRoutes {
	return AuthRoutes{
		logger:         logger,
		handler:        handler,
		authController: authController,
	}
}

func (u AuthRoutes) Setup() {
	u.logger.Info("Setting auth routes")
	{
		u.handler.Gin.POST("/login", u.authController.Login)
		//u.handler.Gin.POST("/register", u.authController.Register)
	}
}
