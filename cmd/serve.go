package cmd

import (
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/api/routes"
	"github.com/721945/dlaw-backend/libs"
)

func RunInit(
	env libs.Env,
	router libs.RequestHandler,
	logger *libs.Logger,
	route routes.Routes,
	middlewares middlewares.Middlewares,
) {

	middlewares.Setup()
	route.Setup()

	logger.Info("🚀 Server is running on port " + env.ServerPort)

	_ = (*router.Gin).Run(":" + env.ServerPort)
}

//func RunMigration(db libs.Database, logger*libs.Logger) {
//	logger.Info("Start Migrations")
//	db.DB.AutoMigrate(&models.User{})
//	logger.Info("Migrations ran successfully")
//
//}
