package natsrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventsv1 "short-link/proto/gen/events/v1"

	"github.com/ananaslegend/short-link/internal/app/config"
	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type Repository struct {
	js     nats.JetStreamContext
	cfg    config.Config
	logger zerolog.Logger
}

func NewRepository(js nats.JetStreamContext, cfg config.Config, logger zerolog.Logger) *Repository {
	return &Repository{js: js, cfg: cfg, logger: logger}
}

func (r Repository) AddRedirectEvent(ctx context.Context, e domain.RedirectEventStatistic) error {
	eventID := uuid.NewString()

	payload := &eventsv1.RedirectEventV1{
		EventId:    eventID,
		OccurredAt: timestamppb.Now(),
		Alias:      e.Alias,
		LinkId:     0, // TODO: should come from domain if needed
		Link:       e.Link,
	}

	body, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = r.js.Publish(
		"events.redirect.v1",
		body,
		nats.AckWait(r.cfg.Nats.PubAckWait),
		nats.MsgId(eventID),
	)

	return err
}
