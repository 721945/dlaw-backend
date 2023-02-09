package middlewares

import (
	"github.com/721945/dlaw-backend/constants"
	"github.com/721945/dlaw-backend/libs"
	"net/http"

	//"github.com/dipeshdulal/clean-gin/constants"
	//"github.com/dipeshdulal/clean-gin/lib"
	"github.com/gin-gonic/gin"
)

// DatabaseTrx middleware for transactions support for database
type DatabaseTrx struct {
	handler libs.RequestHandler
	logger  libs.Logger
	db      libs.Database
}

// statusInList function checks if context writer status is in provided list
func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// NewDatabaseTrx creates new database transactions middleware
func NewDatabaseTrx(
	handler libs.RequestHandler,
	logger libs.Logger,
	db libs.Database,
) DatabaseTrx {
	return DatabaseTrx{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}

// Setup sets up database transaction middleware
func (m DatabaseTrx) Setup() {
	//m.logger.Info("setting up database transaction middleware")

	m.handler.Gin.Use(func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		//m.logger.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		// commit transaction on success status
		if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("Committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("Trx commit error: ", err)
			}
		} else {
			m.logger.Info("Rolling back transaction due to status code: 500")
			txHandle.Rollback()
		}
	})
}
