package handler

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	natswrapper "github.com/ananaslegend/statistic/internal/app/nats_wrapper"
	"github.com/ananaslegend/statistic/internal/statistic/domain"
	"github.com/ananaslegend/statistic/internal/statistic/service"
	eventsv1 "short-link/proto/gen/events/v1"
)

func Module() fx.Option {
	return fx.Module(
		"internal.statistic.handler",
		fx.Provide(NewConsumer),
	)
}

func NewConsumer(redirectHandler service.RedirectHandler, errorHandler natswrapper.JetStreamErrorHandler) *Consumer {
	return &Consumer{
		redirectHandler: redirectHandler,
		errorHandler:    errorHandler,
	}
}

const (
	maxEventDelivers = 3
)

type Consumer struct {
	redirectHandler service.RedirectHandler

	errorHandler natswrapper.JetStreamErrorHandler
}

func (c Consumer) HandleMsg(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		var e eventsv1.RedirectEventV1
		if err := proto.Unmarshal(msg.Data(), &e); err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("failed to unmarshal redirect event")

			c.errorHandler.HandleJetStreamError(ctx, msg, fmt.Errorf("failed to unmarshal redirect event: %v", err))

			return
		}

		err := c.redirectHandler.AddRedirectEvent(ctx, domain.RedirectEventStatistic{
			LinkID:     e.GetLinkId(),
			Link:       e.GetLink(),
			OccurredAt: e.GetOccurredAt().AsTime(),
			Alias:      e.GetAlias(),
		})
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("failed to handle redirect event")
			c.errorHandler.HandleJetStreamError(ctx, msg, fmt.Errorf("failed to handle redirect event: %v", err))

			return
		}

		if err = msg.Ack(); err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("failed to handle ack redirect event")
		}
	}
}

