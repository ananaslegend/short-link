package redirect

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/internal/statistic"
)

type SelectLinkRepo interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type StatManager interface {
	AppendRow(row *statistic.Row)
}

type UseCase struct {
	repo        SelectLinkRepo
	statManager StatManager
}

func NewUseCase(lp SelectLinkRepo, stat StatManager) *UseCase {
	return &UseCase{
		repo:        lp,
		statManager: stat,
	}
}

func (uc UseCase) GetLink(ctx context.Context, alias string) (string, error) {
	const op = "services.link.GetLink"

	var (
		rowStat = statistic.NewRow()
	)

	if alias == "" {
		return "", ErrEmptyAlias
	}

	link, err := uc.repo.SelectLink(ctx, alias)
	if err != nil {
		if errors.Is(err, ErrCantSetToCache) {
			return link, fmt.Errorf("%s: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	rowStat.Link = link
	rowStat.Redirect += 1
	uc.statManager.AppendRow(rowStat)

	return link, nil
}
