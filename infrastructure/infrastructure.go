package infrastructure

import (
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(google_vision.NewGoogleVision),
	fx.Provide(google_storage.NewGoogleStorage),
)
