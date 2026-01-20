package nats_wrapper

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/config"
)

type JetStream nats.JetStreamContext

func Module() fx.Option {
	return fx.Module(
		"internal.app.nats",
		fx.Provide(newNatsConn),
		fx.Provide(newJetStream),
		fx.Invoke(ensureStreams),
	)
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

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			nc.Drain()
			nc.Close()

			return nil
		},
	})

	return nc
}

func newJetStream(
	lc fx.Lifecycle,
	cfg config.Config,
	logger zerolog.Logger,
	nc *nats.Conn,
) nats.JetStreamContext {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get JetStream context")
	}

	return js
}

func ensureStreams(logger zerolog.Logger, js nats.JetStreamContext) {
	streams := []nats.StreamConfig{
		{
			Name:       "events",
			Subjects:   []string{"events.>"},
			Retention:  nats.WorkQueuePolicy,
			Storage:    nats.FileStorage,
			Discard:    nats.DiscardOld,
			Replicas:   1,
			Duplicates: time.Hour,
		},
		{
			Name:       "DLQ",
			Subjects:   []string{"DLQ.>"},
			Retention:  nats.WorkQueuePolicy,
			Storage:    nats.FileStorage,
			Discard:    nats.DiscardOld,
			Replicas:   1,
			Duplicates: time.Hour,
		},
	}

	for _, cfg := range streams {
		_, err := js.AddStream(&cfg)
		if err != nil {
			// Stream might already exist, try to update
			_, err = js.UpdateStream(&cfg)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("stream", cfg.Name).
					Msg("failed to create or update stream")
			} else {
				logger.Info().Str("stream", cfg.Name).Msg("stream updated")
			}
		} else {
			logger.Info().Str("stream", cfg.Name).Msg("stream created")
		}
	}
}
