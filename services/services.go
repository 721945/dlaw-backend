package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewJWTAuthService),
	fx.Provide(NewActionService),
	fx.Provide(NewAppointmentService),
	fx.Provide(NewCaseService),
	fx.Provide(NewFileService),
	fx.Provide(NewFolderService),
	fx.Provide(NewPermissionService),
	fx.Provide(NewTagService),
)
