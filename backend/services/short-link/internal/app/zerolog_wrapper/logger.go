package zerolog_wrapper

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func SetupZerolog(cfg config.Config) zerolog.Logger {
	logger := zerolog.New(zerolog.MultiLevelWriter(
		LevelWriter{
			Levels: []zerolog.Level{zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel},
			Writer: os.Stderr,
		},
		LevelWriter{
			Levels: []zerolog.Level{
				zerolog.WarnLevel,
				zerolog.InfoLevel,
				zerolog.DebugLevel,
				zerolog.TraceLevel,
			},
			Writer: os.Stdout,
		},
	))

	if cfg.Environment == config.Local {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return logger.Level(logLevel(cfg)).
		With().Timestamp().Logger().
		With().Stack().Logger().
		With().Str("service", cfg.ServiceName).Logger().
		With().Str("env", string(cfg.Environment)).Logger()
}

func logLevel(cfg config.Config) zerolog.Level {
	switch cfg.Environment {
	case config.Dev:
		return zerolog.InfoLevel

	case config.Local:
		return zerolog.TraceLevel

	default:
		return zerolog.ErrorLevel
	}
}

type LevelWriter struct {
	Levels []zerolog.Level
	io.Writer
}

func (w LevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	for _, l := range w.Levels {
		if l == level {
			return w.Write(p)
		}
	}

	return len(p), nil
}
