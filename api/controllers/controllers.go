package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserController),
	fx.Provide(NewAuthController),
	fx.Provide(NewActionController),
	fx.Provide(NewFileController),
	fx.Provide(NewTagController),
	fx.Provide(NewPermissionController),
	fx.Provide(NewCaseController),
)
