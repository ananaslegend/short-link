package http

import (
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type LinkHandler struct {
	linkGetter   LinkGetter
	linkInserter LinkInserter
	tracer       trace.Tracer
}

func NewHandler(
	linkGetter LinkGetter,
	linkInserter LinkInserter,
	traceProvider *sdktrace.TracerProvider,
) *LinkHandler {
	return &LinkHandler{
		linkGetter:   linkGetter,
		linkInserter: linkInserter,
		tracer:       traceProvider.Tracer("internal.link.handler.http.LinkHandler"),
	}
}
