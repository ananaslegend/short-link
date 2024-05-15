package service

import (
	"context"
	"github.com/ananaslegend/short-link/internal/statistic"
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

func New(log *slog.Logger, lp SelectLinkRepo, stat StatManager) *Service {
	return &Service{
		log:         log,
		repo:        lp,
		statManager: stat,
	}
}

func (s Service) GetLink(ctx context.Context, alias string) (string, error) {
	link, err := s.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", err
	}

	s.addStatistic(ctx, link)

	return link, nil
}

func (s Service) addStatistic(ctx context.Context, link string) {
	rowStat := statistic.NewRow()
	rowStat.Link = link
	rowStat.Redirect += 1
	s.statManager.AppendRow(rowStat)
}
