package libs

import (
	"fmt"
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
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

	//engine.Use(gin.Recovery())
	engine.Use(gin.CustomRecovery(errorHandler))
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

func errorHandler(c *gin.Context, err any) {
	// check if err is error type
	if _, ok := err.(error); ok {
		// do something
		c.AbortWithStatusJSON(controllers.StatusCode(err.(error)), gin.H{"error": err.(error).Error()})
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
