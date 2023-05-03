package middlewares

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/gin-contrib/cors"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler libs.RequestHandler
	logger  *libs.Logger
	env     libs.Env
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler libs.RequestHandler, logger *libs.Logger, env libs.Env) CorsMiddleware {
	return CorsMiddleware{
		handler: handler,
		logger:  logger,
		env:     env,
	}
}

// Setup sets up cors middleware
func (m CorsMiddleware) Setup() {
	m.logger.Info("Setting up cors middleware")

	debug := m.env.Environment == "development"

	if debug {
		m.logger.Info("Cors is enabled in development mode")
	}

	m.handler.Gin.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "HEAD", "DELETE", "OPTIONS"},
	}))
}
