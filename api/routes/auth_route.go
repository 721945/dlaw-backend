package routes

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/libs"
)

type AuthRoute struct {
	logger         *libs.Logger
	handler        libs.RequestHandler
	authController controllers.AuthController
	authMiddleware middlewares.JWTAuthMiddleware
}

func NewAuthRoute(
	logger *libs.Logger,
	handler libs.RequestHandler,
	authController controllers.AuthController,
	authMiddleware middlewares.JWTAuthMiddleware,
) AuthRoute {
	return AuthRoute{
		logger:         logger,
		handler:        handler,
		authController: authController,
		authMiddleware: authMiddleware,
	}
}

func (u AuthRoute) Setup() {
	u.logger.Info("Setting auth routes")
	{
		u.handler.Gin.POST("/login", u.authController.Login)
		u.handler.Gin.POST("/forget-password", u.authController.ForgetPassword)
		u.handler.Gin.POST("/reset-password", u.authController.ResetPassword)
		u.handler.Gin.POST("/change-password", u.authController.ChangePassword).Use(u.authMiddleware.Handler())
	}
}
