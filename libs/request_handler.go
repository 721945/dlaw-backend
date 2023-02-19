package libs

import (
	"fmt"
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
