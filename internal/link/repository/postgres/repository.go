package postgres

import (
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	pgxwrapper "github.com/ananaslegend/short-link/internal/app/pgx_wrapper"
)

type LinkRepository struct {
	db     *pgxwrapper.Wrapper
	tracer trace.Tracer
}

func NewLinkRepository(
	db *pgxwrapper.Wrapper,
	traceProvider *sdktrace.TracerProvider,
) *LinkRepository {
	return &LinkRepository{
		db:     db,
		tracer: traceProvider.Tracer("internal.link.repository.postgres.LinkRepository"),
	}
}
