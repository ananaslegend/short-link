package transport

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	natswrapper "github.com/ananaslegend/statistic/internal/app/nats_wrapper"
	"github.com/ananaslegend/statistic/internal/statistic/handler"
)

func Module() fx.Option {
	return fx.Module(
		"internal.transport",
		fx.Invoke(RegisterRedirectConsumer),
	)
}

func RegisterRedirectConsumer(ctx context.Context, logger zerolog.Logger, js jetstream.JetStream, consumer *handler.Consumer) {
	c := natswrapper.NewConsumer(ctx, logger, js, "events", jetstream.ConsumerConfig{
		Durable:       "statistic_redirect_consumer_v1",
		FilterSubject: "events.redirect.v1",
	})

	go c.Consume(consumer.HandleMsg(ctx))
}
