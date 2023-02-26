package routes

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewActionRoute),
	fx.Provide(NewAppointmentRoute),
	fx.Provide(NewAuthRoute),
	fx.Provide(NewCaseRoute),
	fx.Provide(NewFileRoute),
	fx.Provide(NewFileTypeRoute),
	fx.Provide(NewFolderRoute),
	fx.Provide(NewPermissionRoute),
	fx.Provide(NewTagRoute),
	fx.Provide(NewUserRoute),
	fx.Provide(NewRoutes),
)

type Routes []Route

type Route interface {
	Setup()
}

func NewRoutes(
	actionRoute ActionRoute,
	appointmentRoute AppointmentRoute,
	authRoute AuthRoute,
	caseRoute CaseRoute,
	fileRoute FileRoute,
	fileTypeRoute FileTypeRoute,
	folderRoute FolderRoute,
	permissionRoute PermissionRoute,
	tagRoute TagRoute,
	userRoute UserRoute,
) Routes {
	return Routes{
		actionRoute,
		appointmentRoute,
		authRoute,
		caseRoute,
		fileRoute,
		fileTypeRoute,
		folderRoute,
		permissionRoute,
		tagRoute,
		userRoute,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
