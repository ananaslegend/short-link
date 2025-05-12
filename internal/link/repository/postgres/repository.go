package postgres

import (
	pgxwrapper "github.com/ananaslegend/short-link/internal/app/pgx_wrapper"
)

type LinkRepository struct {
	db *pgxwrapper.Wrapper
}

func NewLinkRepository(
	db *pgxwrapper.Wrapper,
) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}
