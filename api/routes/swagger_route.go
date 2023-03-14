package routes

import (
	_ "github.com/721945/dlaw-backend/docs"
	"github.com/721945/dlaw-backend/libs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRoute struct {
	handler libs.RequestHandler
	logger  *libs.Logger
}

func NewSwaggerRoute(handler libs.RequestHandler, logger *libs.Logger) SwaggerRoute {
	return SwaggerRoute{handler: handler, logger: logger}
}

func (r SwaggerRoute) Setup() {
	r.logger.Info("Setting tag routes")
	r.handler.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
