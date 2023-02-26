package libs

import (
	"fmt"
	"github.com/JosephWoodward/gin-errorhandling/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

type RequestHandler struct {
	Gin    *gin.Engine
	logger *Logger
}

func NewRequestHandler(logger *Logger) RequestHandler {
	//gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()

	gin.ForceConsoleColor()
	logger.Info("Setting up request handler")
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler(
		middleware.Map(ErrInternalServerError).ToResponse(func(c *gin.Context, err error) {
			logger.Error(err)
			c.AbortWithStatusJSON(StatusCode(err.(error)), gin.H{"error": err.Error()})
			return
		}),
	))

	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	return RequestHandler{
		Gin:    engine,
		logger: logger,
	}
}
