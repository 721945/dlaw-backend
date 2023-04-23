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
	newApi := u.handler.Gin.Group("/users")
	{
		newApi.POST("", u.userController.CreateUser)
		newApi.GET("", u.userController.GetUsers)
		newApi.GET("/:id", u.userController.GetUser)
		newApi.GET("/me", u.userController.GetMe).Use(u.authMiddleware.Handler())
		newApi.PATCH("/:id", u.userController.UpdateUser)
		newApi.PATCH("/me", u.userController.UpdateUser).Use(u.authMiddleware.Handler())

	}
	//}	api := u.handler.Gin.Group("/users")
	//{
	//	api.POST("", u.userController.CreateUser)
	//	api.GET("", u.userController.GetUsers)
	//	api.GET("/:id", u.userController.GetUser)
	//	api.GET("/me", u.userController.GetMe).Use(u.authMiddleware.Handler())
	//	api.PUT("/:id", u.userController.UpdateUser)
	//}
	//
}
