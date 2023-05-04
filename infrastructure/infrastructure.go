package infrastructure

import (
	"github.com/721945/dlaw-backend/infrastructure/google_calendar"
	"github.com/721945/dlaw-backend/infrastructure/google_storage"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"github.com/721945/dlaw-backend/infrastructure/smtp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(google_vision.NewGoogleVision),
	fx.Provide(google_storage.NewGoogleStorage),
	fx.Provide(smtp.NewSMTP),
	fx.Provide(google_calendar.NewGoogleCalendar),
)
