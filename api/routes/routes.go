package routes

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Provide(NewRoutes),
)

type Routes []Route

type Route interface {
	Setup()
}

func NewRoutes(
	userRoute UserRoutes,
) Routes {
	return Routes{
		userRoute,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
