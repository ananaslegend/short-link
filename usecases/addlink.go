package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/core"
	"github.com/ananaslegend/short-link/errs"
	"github.com/ananaslegend/short-link/logs"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"strings"
)

type LinkSaver interface {
	SaveLink(link, alias string) error
}

func AddLink(c context.Context, log *slog.Logger, ls LinkSaver, link, alias string) (string, error) {
	const op = "usecases.link.SaveLink"

	var autoAlias bool
	if len(alias) == 0 {
		alias = core.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	if err := ls.SaveLink(link, alias); err != nil {
		if errors.Is(err, errs.ErrAliasExists) {
			if autoAlias {
				log.Error(fmt.Sprintf("auto generated alias already exists"), logs.Err(err))
				return "", errs.ErrAutoAliasAlreadyExists
			}

			err = fmt.Errorf("%s: %w", op, err)
		}
		return "", err
	}

	return alias, nil
}
