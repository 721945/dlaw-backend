package repositories

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewActionRepository),
	fx.Provide(NewAppointmentRepository),
	fx.Provide(NewCaseRepository),
	fx.Provide(NewFileRepository),
	fx.Provide(NewFileTypeRepository),
	fx.Provide(NewFileUrlRepository),
	fx.Provide(NewFolderRepository),
	fx.Provide(NewPermissionRepository),
	fx.Provide(NewTagRepository),
	fx.Provide(NewCasePermissionLogRepository),
	fx.Provide(NewCasePermissionRepository),
	fx.Provide(NewActionLogRepository),
	fx.Provide(NewCaseUsedLogRepository),
)
