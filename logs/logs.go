package logs

import (
	"github.com/ananaslegend/short-link/config"
	"log/slog"
	"os"
)

func SetUpLogger(cfg config.AppConfig) *slog.Logger {
	var logger *slog.Logger

	switch cfg.Env {
	case config.Local:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
