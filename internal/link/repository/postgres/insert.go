package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ananaslegend/short-link/internal/link/domain"
)

func (r LinkRepository) InsertAliasedLink(
	ctx context.Context,
	dto domain.InsertAliasedLink,
) (domain.AliasedLink, error) {
	const op = "internal.link.postgres.Repository.InsertLink"

	res := domain.AliasedLink{}

	err := r.db.Pool.QueryRow(ctx,
		`insert into link (alias, link) values ($1,$2) returning id, alias, link;`,
		dto.Alias, dto.Link).
		Scan(&res.ID, &res.Alias, &res.Link)
	if err != nil {
		if pgxErr := (&pgconn.PgError{}); errors.As(err, &pgxErr) &&
			pgxErr.Code == pgerrcode.UniqueViolation {
			return res, ErrAliasNtFound
		}

		return domain.AliasedLink{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
