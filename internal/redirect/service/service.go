package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/internal/redirect/repository"
	"github.com/ananaslegend/short-link/internal/statistic"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
)

type SelectLinkRepo interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type StatManager interface {
	AppendRow(row *statistic.Row)
}

type Service struct {
	log         *slog.Logger
	repo        SelectLinkRepo
	statManager StatManager
}

func NewService(lp SelectLinkRepo, stat StatManager) *Service {
	return &Service{
		repo:        lp,
		statManager: stat,
	}
}

func (uc Service) GetLink(ctx context.Context, alias string) (string, error) {
	const op = "internal.redirect.service.Service.GetLink"
	logger := uc.log.With(slog.String("op", op))

	var (
		rowStat = statistic.NewRow()
	)

	link, err := uc.repo.SelectLink(ctx, alias)
	if err != nil {
		if errors.Is(err, repository.ErrCantSetToCache) {
			logger.Error("failed to set link to cache", logs.Err(err))
		} else {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	rowStat.Link = link
	rowStat.Redirect += 1
	uc.statManager.AppendRow(rowStat)

	return link, nil
}
