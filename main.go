/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/721945/dlaw-backend/cmd"
	_ "github.com/721945/dlaw-backend/docs"
)

// @title Swagger Example API
// @version 2.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	cmd.Execute()
	//env := libs.NewEnv()
	//logger := libs.NewLogger()
	//db := libs.NewDatabase(env, logger)
	//db.DB.AutoMigrate(&models.User{})
}
