package cmd

import (
	"github.com/721945/dlaw-backend/api/routes"
	"github.com/721945/dlaw-backend/libs"
)

func RunInit(env libs.Env, router libs.RequestHandler, logger libs.Logger, route routes.Routes) {
	route.Setup()
	println("🚀 Server is running on port " + env.ServerPort)
	router.Gin.Run(":" + env.ServerPort)
}

//func RunMigration(db libs.Database, logger libs.Logger) {
//	logger.Info("Start Migrations")
//	db.DB.AutoMigrate(&models.User{})
//	logger.Info("Migrations ran successfully")
//
//}
