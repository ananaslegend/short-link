package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/ananaslegend/short-link/internal/save/repository"
	"github.com/ananaslegend/short-link/pkg/cslog"
	"github.com/ananaslegend/short-link/pkg/shortner"
)

type InsertLinkRepo interface {
	InsertLink(ctx context.Context, link, alias string) error
}

type Service struct {
	repo InsertLinkRepo
}

func New(lp InsertLinkRepo) *Service {
	return &Service{
		repo: lp,
	}
}

func (s Service) AddLink(ctx context.Context, link, alias string) (string, error) {
	if len(alias) == 0 {
		alias = shortner.MakeShorter(uuid.New().ID())
		ctx = cslog.With(ctx, slog.Bool("auto_generated_alias", true))
	}

	ctx = cslog.With(ctx, slog.String("alias", alias))

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	ctx = cslog.With(ctx, slog.String("link", link))

	if err := s.repo.InsertLink(ctx, link, alias); err != nil {
		cslog.Logger(ctx).Error("insert link", cslog.Error(err))

		if errors.Is(err, repository.ErrAliasAlreadyExists) {
			return "", ErrAliasAlreadyExists
		}

		return "", err
	}

	return alias, nil
}
