package bootstrap

import (
	"github.com/721945/dlaw-backend/api/controllers"
	"github.com/721945/dlaw-backend/api/middlewares"
	"github.com/721945/dlaw-backend/api/routes"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/721945/dlaw-backend/services"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	repositories.Module,
	services.Module,
	controllers.Module,
	libs.Module,
	middlewares.Module,
	routes.Module,
)
