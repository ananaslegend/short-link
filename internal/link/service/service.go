package service

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Link struct {
	linkGetter          LinkGetter
	aliasedLinkInserter AliasedLinkInserter
	aliasGenerator      AliasGenerator

	tracer trace.Tracer
}

func New(
	lp LinkGetter,
	aliasedLinkInserter AliasedLinkInserter,
	aliasGenerator AliasGenerator,
	traceProvider *sdktrace.TracerProvider,
) *Link {
	return &Link{
		linkGetter:          lp,
		aliasedLinkInserter: aliasedLinkInserter,
		aliasGenerator:      aliasGenerator,
		tracer:              traceProvider.Tracer("internal.link.service.Link"),
	}
}
