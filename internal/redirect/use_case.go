package redirect

import (
	"context"
	"fmt"
)

type SelectLinkRepo interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type UseCase struct {
	repo SelectLinkRepo
}

func NewUseCase(lp SelectLinkRepo) *UseCase {
	return &UseCase{
		repo: lp,
	}
}

func (uc UseCase) GetLink(ctx context.Context, alias string) (string, error) {
	const op = "services.link.Get"

	if alias == "" {
		return "", ErrEmptyAlias
	}

	link, err := uc.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
