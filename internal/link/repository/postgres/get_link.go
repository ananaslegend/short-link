package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r LinkRepository) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "internal.link.redirect.repository.postgres.Repository.GetLinkByAlias"

	var link string

	err := r.db.Pool.QueryRow(ctx, `select link from link where alias = $1`, alias).Scan(&link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAliasNtFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
