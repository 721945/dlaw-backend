package middlewares

import (
	"github.com/721945/dlaw-backend/libs"
)

// DatabaseTrx middleware for transactions support for database
type DatabaseTrx struct {
	handler libs.RequestHandler
	logger  *libs.Logger
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
	logger *libs.Logger,
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
	m.logger.Info("1 - setting up database transaction middleware")

	//m.handler.Gin.Use(DBTransactionMiddleware(m.db.DB))
	//m.handler.Gin.Use(
	//	func(c *gin.Context) {
	//	m.logger.Info("2 - setting up database transaction middleware")
	//	txHandle := m.db.DB.Begin()
	//	//m.logger.Info("beginning database transaction")
	//
	//	// rollback transaction on panic
	//	m.logger.Info("Rolling back transaction on panic")
	//	defer func() {
	//		if r := recover(); r != nil {
	//			txHandle.Rollback()
	//		}
	//	}()
	//
	//	// set transaction in context
	//	m.logger.Info("Setting db transaction in context 🔌")
	//	c.Set(constants.DBTransaction, txHandle)
	//	m.logger.Info("Setting db transaction in context 🔌")
	//	c.Next()
	//
	//	// commit transaction on success status
	//	if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
	//		m.logger.Info("Committing transactions")
	//		if err := txHandle.Commit().Error; err != nil {
	//			m.logger.Error("Trx commit error: ", err)
	//		}
	//	} else {
	//		m.logger.Info("Rolling back transaction due to status code: 500")
	//		txHandle.Rollback()
	//	}
	//})
}

//func DBTransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		txHandle := db.Begin()
//		log.Print("beginning database transaction")
//
//		defer func() {
//			if r := recover(); r != nil {
//				txHandle.Rollback()
//			}
//		}()
//
//		c.Set(constants.DBTransaction, txHandle)
//		c.Next()
//
//		//if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
//		//	log.Print("committing transactions")
//		//	if err := txHandle.Commit().Error; err != nil {
//		//		log.Print("trx commit error: ", err)
//		//	}
//		//} else {
//		//	log.Print("rolling back transaction due to status code: ", c.Writer.Status())
//		//	txHandle.Rollback()
//		//}
//	}
//}
