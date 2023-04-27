package middlewares

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/gin-gonic/gin"
)

type PaginationMiddleware struct {
	logger *libs.Logger
}

func NewPaginationMiddleware(
	logger *libs.Logger,
) PaginationMiddleware {
	return PaginationMiddleware{logger: logger}
}

func (m PaginationMiddleware) Setup() {
	m.logger.Info("Setting up jwt auth middleware")
}

func (m PaginationMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagination dtos.Pagination

		if err := c.ShouldBindQuery(&pagination); err != nil {
			c.Next()
			return
		}

		c.Set("pagination", pagination)

		c.Next()
	}
}
