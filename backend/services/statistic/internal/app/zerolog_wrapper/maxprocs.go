package zerolog_wrapper

import "github.com/rs/zerolog"

func ZerologMaxProcsLogger(logger zerolog.Logger) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		logger.Info().Str("component", "automaxprocs").Msgf(format, args...)
	}
}
