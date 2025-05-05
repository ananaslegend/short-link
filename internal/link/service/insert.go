package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	aliasdomain "github.com/ananaslegend/short-link/internal/alias_generator/domain"
	"github.com/ananaslegend/short-link/internal/link/domain"
	"github.com/ananaslegend/short-link/internal/link/repository/postgres"
)

type AliasGenerator interface {
	GenerateAlias(context context.Context, alias aliasdomain.GenerateAlias) (string, error)
}

type AliasedLinkInserter interface {
	InsertAliasedLink(ctx context.Context, dto domain.InsertAliasedLink) (domain.AliasedLink, error)
}

func (s Link) InsertLink(ctx context.Context, dto domain.InsertLink) (domain.AliasedLink, error) {
	const op = "short-link.internal.link.service.AliasedLink.InsertLink"

	if dto.Alias == nil {
		if err := s.generateAlias(ctx, &dto); err != nil {
			zerolog.Ctx(ctx).
				Error().
				Err(err).
				Str("op", op).
				Any("input", dto).
				Msg("failed to generate alias")

			return domain.AliasedLink{}, fmt.Errorf("%v: %w", op, err)
		}
	}

	insertedLink, err := s.aliasedLinkInserter.InsertAliasedLink(ctx, mapToInsertAliasedLink(dto))
	if err != nil {
		return domain.AliasedLink{}, handleInsertionError(ctx, err, dto)
	}

	return insertedLink, nil
}

func (s Link) generateAlias(ctx context.Context, dto *domain.InsertLink) error {
	alias, err := s.aliasGenerator.GenerateAlias(ctx, aliasdomain.GenerateAlias{
		Link: dto.Link,
	})
	if err != nil {
		return err
	}

	dto.Alias = &alias

	return nil
}

func mapToInsertAliasedLink(dto domain.InsertLink) domain.InsertAliasedLink {
	return domain.InsertAliasedLink{
		Link:  dto.Link,
		Alias: *dto.Alias,
	}
}

func handleInsertionError(ctx context.Context, err error, dto domain.InsertLink) error {
	const op = "short-link.internal.link.service.handleInsertionError"

	if errors.Is(err, postgres.ErrAliasAlreadyExists) {
		if dto.Alias == nil {
			zerolog.Ctx(ctx).
				Error().
				Err(err).
				Str("op", op).
				Any("input", dto).
				Msg("failed to insert aliased link")
		}

		return ErrAliasAlreadyExists
	}

	zerolog.Ctx(ctx).
		Error().
		Err(err).
		Str("op", op).
		Any("input", dto).
		Msg("failed to insert aliased link")

	return fmt.Errorf("%v: %w", op, err)
}

func (s Link) InsertAliasedLink(
	ctx context.Context,
	dto domain.InsertAliasedLink,
) (domain.AliasedLink, error) {
	const op = "short-link.internal.link.service.AliasedLink.InsertAliasedLink"

	link, err := s.aliasedLinkInserter.InsertAliasedLink(ctx, dto)
	if err != nil {
		zerolog.Ctx(ctx).
			Error().
			Err(err).
			Str("op", op).
			Any("input", dto).
			Msg("failed to insert link")

		return domain.AliasedLink{}, fmt.Errorf("%v: %w", op, err)
	}

	return link, nil
}
