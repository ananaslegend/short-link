package zerolog_wrapper

import "github.com/rs/zerolog"

type ZerologLevelWriter struct {
	Logger zerolog.Logger
	Level  zerolog.Level
}

// Write implements the io.Writer interface for zerolog with Level instead of zerolog.NoLevel like original.
func (l ZerologLevelWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
	}

	l.Logger.WithLevel(l.Level).Msg(string(p))

	return n, nil
}
