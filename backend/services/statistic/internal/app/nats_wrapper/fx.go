package nats_wrapper

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"time"

	"github.com/ananaslegend/statistic/internal/app/config"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type JetStream nats.JetStreamContext

func Module() fx.Option {
	return fx.Module(
		"internal.app.nats",
		fx.Provide(newNatsConn),
		fx.Provide(newJetStream),
		fx.Provide(
			fx.Annotate(NewDLQPublisher, fx.As(new(JetStreamErrorHandler))),
		),
	)
}

// JetStreamErrorHandler interface for DLQ handling
type JetStreamErrorHandler interface {
	HandleJetStreamError(ctx context.Context, msg jetstream.Msg, handleErr error)
}

func newNatsConn(lc fx.Lifecycle, cfg config.Config, logger zerolog.Logger) *nats.Conn {
	opts := []nats.Option{
		nats.Name(cfg.ServiceName),
		nats.Timeout(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
	}

	nc, err := nats.Connect(cfg.Nats.URL, opts...)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to NATS")
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				nc.Drain()
				nc.Close()

				return nil
			},
		},
	)

	return nc
}

func newJetStream(logger zerolog.Logger, nc *nats.Conn) jetstream.JetStream {
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get JetStream")
	}

	return js
}

func NewConsumer(ctx context.Context, logger zerolog.Logger, js jetstream.JetStream, stream string, cfg jetstream.ConsumerConfig) jetstream.Consumer {
	c, err := js.CreateOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		logger.Fatal().Err(err).Str("stream", stream).Msg("failed to create or update JetStream consumer")
	}

	return c
}
