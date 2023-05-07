package libs

import (
	_ "github.com/721945/dlaw-backend/docs"
	"github.com/JosephWoodward/gin-errorhandling/middleware"
	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Gin    *gin.Engine
	logger *Logger
}

func NewRequestHandler(logger *Logger) RequestHandler {
	logger.Info("Setting up request handler")

	engine := gin.Default()

	gin.ForceConsoleColor()

	//engine.Use(gin.Recovery())

	engine.Use(middleware.ErrorHandler(
		middleware.Map(ErrInternalServerError).ToResponse(func(c *gin.Context, err error) {
			logger.Error(err)
			c.AbortWithStatusJSON(StatusCode(err.(error)), gin.H{"error": err.Error()})
			return
		}),
	))
	//

	return RequestHandler{
		Gin:    engine,
		logger: logger,
	}
}
