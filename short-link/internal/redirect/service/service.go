package service

import (
	"context"
	"github.com/ananaslegend/short-link/internal/statistic"
)

type SelectLinkRepo interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type StatManager interface {
	AppendRow(row *statistic.Row)
}

type Service struct {
	repo        SelectLinkRepo
	statManager StatManager
}

func New(lp SelectLinkRepo, stat StatManager) *Service {
	return &Service{
		repo:        lp,
		statManager: stat,
	}
}

func (s Service) GetLink(ctx context.Context, alias string) (string, error) {
	link, err := s.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", err
	}

	addStatistic(ctx, link, alias)

	return link, nil
}

func addStatistic(ctx context.Context, link string, alias string) {
	if rowStat, ok := statistic.GetFromCtx(ctx); ok {
		rowStat.Link = link
		rowStat.Alias = alias
		rowStat.Redirect += 1
	}
}
