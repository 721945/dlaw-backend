package middlewares

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JWTAuthMiddleware struct {
	service services.JWTAuthService
	logger  *libs.Logger
}

func NewJWTAuthMiddleware(
	service services.JWTAuthService,
	logger *libs.Logger,
) JWTAuthMiddleware {
	return JWTAuthMiddleware{service: service, logger: logger}
}

func (m JWTAuthMiddleware) Setup() {
	m.logger.Info("Setting up jwt auth middleware")
}

func (m JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]

			user, err := m.service.VerifyToken(authToken)

			if err != nil {
				_ = c.Error(libs.ErrUnauthorized)
				c.Abort()
				return
			}

			m.logger.Info("USER : ", user.ID)

			c.Set("user", user)
			c.Set("id", user.ID)
			c.Next()

			return

		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized !!",
		})

		c.Abort()
	}
}
