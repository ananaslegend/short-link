package nats_wrapper

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"slices"
)

const (
	defaultMaxEventDelivers = 3
)

type DLQMsg struct {
	OriginalSubject string
	Data            []byte
	Error           string
}

type DLQPublisher struct {
	dlqJS jetstream.JetStream

	retryableErrors  []error
	maxEventDelivers uint64
}

type DLQPublisherConfig struct {
	retryableErrors  []error
	maxEventDelivers uint64
}

// NewDLQPublisher creates a new DLQPublisher with the provided JetStream context and optional configuration.
// Provide only 1 configuration struct.
func NewDLQPublisher(js jetstream.JetStream, cfg ...DLQPublisherConfig) *DLQPublisher {
	dql := &DLQPublisher{
		dlqJS:            js,
		maxEventDelivers: defaultMaxEventDelivers,
	}

	if len(cfg) > 0 {
		applyConfig(dql, cfg[0])
	}

	return dql
}

func applyConfig(dql *DLQPublisher, cfg DLQPublisherConfig) {
	if cfg.retryableErrors != nil {
		dql.retryableErrors = cfg.retryableErrors
	}

	cfg.maxEventDelivers = defaultMaxEventDelivers

	if cfg.maxEventDelivers != 0 {
		dql.maxEventDelivers = cfg.maxEventDelivers
	}
}

func (p DLQPublisher) HandleJetStreamError(ctx context.Context, msg jetstream.Msg, handleErr error) {
	if slices.Contains(p.retryableErrors, handleErr) {
		metadata, err := msg.Metadata()
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("failed to get metadata for redirect event")

			p.publishDQL(ctx, msg, handleErr)

			return
		}

		if metadata.NumDelivered <= p.maxEventDelivers {
			if err = msg.Nak(); err != nil {
				zerolog.Ctx(ctx).Err(handleErr).Msg("failed to handle nak redirect event")
			}
		}
	}

	p.publishDQL(ctx, msg, handleErr)
}

func (p DLQPublisher) publishDQL(ctx context.Context, msg jetstream.Msg, handleErr error) {
	dlqMsg := &DLQMsg{
		OriginalSubject: msg.Subject(),
		Data:            msg.Data(),
		Error:           handleErr.Error(),
	}

	jsonMsg, err := json.Marshal(dlqMsg)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("dlq_message", msg).Msg("failed to marshal dlq message")

		return
	}

	_, err = p.dlqJS.Publish(ctx, "DLQ."+dlqMsg.OriginalSubject, jsonMsg)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("dlq_message", msg).Msg("failed to publish dlq message")

		return
	}

	if err = msg.Ack(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to ack msg after adding to dlq message")
	}
}
