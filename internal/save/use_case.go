package save

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/pkg/shortner"
	"github.com/google/uuid"
	"strings"
)

type InsertLinkRepo interface {
	InsertLink(ctx context.Context, link, alias string) error
}

type UseCase struct {
	repo InsertLinkRepo
}

func NewUseCase(lp InsertLinkRepo) *UseCase {
	return &UseCase{
		repo: lp,
	}
}

func (ls UseCase) AddLink(ctx context.Context, link, alias string) (string, error) {
	const op = "services.link.Add"

	var autoAlias bool
	if len(alias) == 0 {
		alias = shortner.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	if err := ls.repo.InsertLink(ctx, link, alias); err != nil {
		if errors.Is(err, ErrAliasAlreadyExists) {
			if autoAlias {
				return "", fmt.Errorf("%s: %w, alias: %s", op, ErrAutoAliasAlreadyExists, alias)
			}

			err = fmt.Errorf("%s: %w", op, err)
		}
		return "", err
	}

	return alias, nil
}
