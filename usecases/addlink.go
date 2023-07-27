package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/core"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/storage"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type LinkSaver interface {
	AddLink(link, alias string) error
}

func AddLink(c context.Context, log *slog.Logger, ls LinkSaver, link, alias string) (string, error) {
	const op = "usecases.link.AddLink"

	var autoAlias bool
	if len(alias) == 0 {
		alias = core.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if err := ls.AddLink(link, alias); err != nil {
		if errors.Is(err, storage.ErrAliasExists) {
			if !autoAlias {
				log.Error(fmt.Sprintf("auto generated alias already exists"), logs.Err(err))
				return "", ErrAutoAliasAlreadyExists
			}

			err = fmt.Errorf("%s: %w", op, err)
		}
		return "", err
	}

	return alias, nil
}

var ErrAutoAliasAlreadyExists = errors.New("auto alias already exists")
