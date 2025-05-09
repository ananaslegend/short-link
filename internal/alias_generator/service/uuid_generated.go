package service

import (
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"strings"

	"github.com/google/uuid"

	"github.com/ananaslegend/short-link/internal/alias_generator/domain"
)

const alphabet = "1qazQAZ2wsxWSX3edcEDC4rfvRFV5tgbTGB6yhnYHN7ujmUJM8ikK9oLPp"

var alphabetLen = uint32(len(alphabet)) //nolint:gochecknoglobals

type UUIDGenerated struct {
	tracer trace.Tracer
}

func NewUUIDGenerated(traceProvider *sdktrace.TracerProvider) *UUIDGenerated {
	return &UUIDGenerated{tracer: traceProvider.Tracer("internal.alias_generator.UUIDGenerated")}
}

func (u UUIDGenerated) GenerateAlias(
	ctx context.Context,
	alias domain.GenerateAlias,
) (string, error) {
	const op = "internal.alias_generator.UUIDGenerated.GenerateAlias"

	_, span := u.tracer.Start(ctx, op)
	defer span.End()

	return makeShorter(uuid.New().ID()), nil
}

func makeShorter(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
