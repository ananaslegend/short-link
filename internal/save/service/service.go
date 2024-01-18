package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/internal/save/repository"
	"github.com/ananaslegend/short-link/pkg/shortner"
	"github.com/google/uuid"
	"log/slog"
	"strings"
)

type InsertLinkRepo interface {
	InsertLink(ctx context.Context, link, alias string) error
}

type Service struct {
	log  *slog.Logger
	repo InsertLinkRepo
}

func New(log *slog.Logger, lp InsertLinkRepo) *Service {
	return &Service{
		log:  log,
		repo: lp,
	}
}

func (s Service) AddLink(ctx context.Context, link, alias string) (string, error) {
	const op = "services.link.Add"

	var autoAlias bool
	if len(alias) == 0 {
		alias = shortner.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	if err := s.repo.InsertLink(ctx, link, alias); err != nil {
		if errors.Is(err, repository.ErrAliasAlreadyExists) {
			if autoAlias {
				return "", fmt.Errorf("%w, alias: %s", ErrAutoAliasAlreadyExists, alias)
			}

			return "", ErrAliasAlreadyExists
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return alias, nil
}
