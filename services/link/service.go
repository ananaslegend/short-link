package link

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/errs"
	"github.com/ananaslegend/short-link/shortner"
	"github.com/google/uuid"
	"strings"
)

type LinkRepo interface {
	InsertLink(link, alias string) error
	SelectLink(ctx context.Context, alias string) (string, error)
}

type LinkService struct {
	repo LinkRepo
}

func New(lp LinkRepo) *LinkService {
	return &LinkService{
		repo: lp,
	}
}

func (ls LinkService) AddLink(ctx context.Context, link, alias string) (string, error) {
	const op = "services.link.Add"

	var autoAlias bool
	if len(alias) == 0 {
		alias = shortner.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	if err := ls.repo.InsertLink(link, alias); err != nil {
		if errors.Is(err, errs.ErrAliasExists) {
			if autoAlias {
				return "", fmt.Errorf("%s: %w, alias: %s", op, errs.ErrAutoAliasAlreadyExists, alias)
			}

			err = fmt.Errorf("%s: %w", op, err)
		}
		return "", err
	}

	return alias, nil
}

func (ls LinkService) GetLink(ctx context.Context, alias string) (string, error) {
	const op = "services.link.Get"

	if alias == "" {
		return "", errs.ErrEmptyAlias
	}

	link, err := ls.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
