package clog

import (
	"context"
	"log/slog"
)

type Key struct{}

type Clog struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) Clog {
	return Clog{logger: logger}
}

func Init(handler slog.Handler) Clog {
	logger := slog.New(handler)
	return New(logger)
}

func Ctx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(Key{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func With(ctx context.Context, args ...any) context.Context {
	if logger, ok := ctx.Value(Key{}).(*slog.Logger); ok {
		newLogger := logger.With(args) //nolint:govet // todo
		return context.WithValue(ctx, Key{}, newLogger)
	}

	return ctx
}

func WithString(ctx context.Context, key, value string) context.Context {
	if logger, ok := ctx.Value(Key{}).(*slog.Logger); ok {
		newLogger := logger.With(slog.String(key, value))
		return context.WithValue(ctx, Key{}, newLogger)
	}

	return ctx
}

func WithBool(ctx context.Context, key string, value bool) context.Context {
	if logger, ok := ctx.Value(Key{}).(*slog.Logger); ok {
		newLogger := logger.With(slog.Bool(key, value))
		return context.WithValue(ctx, Key{}, newLogger)
	}

	return ctx
}

func WithInt(ctx context.Context, key string, value int) context.Context {
	if logger, ok := ctx.Value(Key{}).(*slog.Logger); ok {
		newLogger := logger.With(slog.Int(key, value))
		return context.WithValue(ctx, Key{}, newLogger)
	}

	return ctx
}

func ErrorMsg(err error) slog.Attr {
	return slog.String("error", err.Error())
}
