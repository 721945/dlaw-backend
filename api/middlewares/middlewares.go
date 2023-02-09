package middlewares

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewDatabaseTrx),
	fx.Provide(NewMiddlewares),
	fx.Provide(NewCorsMiddleware),
)

type IMiddleware interface {
	Setup()
}

type Middlewares []IMiddleware

func NewMiddlewares(
	dbTrxMiddleware DatabaseTrx,
	corsMiddleware CorsMiddleware,
) Middlewares {
	return Middlewares{
		dbTrxMiddleware,
		corsMiddleware,
	}
}

func (m Middlewares) Setup() {
	for _, middleware := range m {
		middleware.Setup()
	}
}
