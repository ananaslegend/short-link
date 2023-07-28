package usecases

import (
	"context"
	"fmt"
	"github.com/ananaslegend/short-link/errs"
	"golang.org/x/exp/slog"
)

type LinkGetter interface {
	GetLink(alias string) (string, error)
}

func GetLink(c context.Context, log *slog.Logger, lg LinkGetter, alias string) (string, error) {
	const op = "usecases.link.SaveLink"

	if alias == "" {
		return "", errs.ErrEmptyAlias
	}

	link, err := lg.GetLink(alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
