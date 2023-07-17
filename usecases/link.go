package usecases

import (
	"context"
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

func AddLink(c context.Context, log *slog.Logger, ls LinkSaver, link, alias string) error {
	const op = "usecases.link.AddLink"
	var autoAlias bool
	// TODO link Validation

	if len(alias) == 0 {
		alias = core.MakeShorter(uuid.New().ID())
		autoAlias = true
	}

	if err := ls.AddLink(link, alias); err != nil {
		if err == storage.ErrAliasExists {
			err = fmt.Errorf("%s: alias %s is already exists", op, alias)
			if autoAlias {
				log.Error("%w", err)
			} // TODO доробити
		}

		err = fmt.Errorf("%s: %w", op, err)
		log.Error("cant add link", logs.Err(err))
		return err
	}

	return nil
}
