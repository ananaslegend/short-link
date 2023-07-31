package link

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/errs"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/shortner"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"strings"
)

type LinkRepo interface {
	InsertLink(link, alias string) error
	SelectLink(alias string) (string, error)
}

type LinkService struct {
	repo LinkRepo
	log  *slog.Logger
}

func New(log *slog.Logger, lp LinkRepo) *LinkService {
	return &LinkService{
		repo: lp,
		log:  log,
	}
}

func (ls LinkService) AddLink(c context.Context, link, alias string) (string, error) {
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
				ls.log.Error(fmt.Sprintf("auto generated alias already exists"), logs.Err(err))
				return "", errs.ErrAutoAliasAlreadyExists
			}

			err = fmt.Errorf("%s: %w", op, err)
		}
		return "", err
	}

	return alias, nil
}

func (ls LinkService) GetLink(c context.Context, alias string) (string, error) {
	const op = "services.link.Get"

	if alias == "" {
		return "", errs.ErrEmptyAlias
	}

	link, err := ls.repo.SelectLink(alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
