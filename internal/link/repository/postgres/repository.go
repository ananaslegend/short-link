package postgres

import (
	pgxwrapper "github.com/ananaslegend/short-link/internal/app/pgx_wrapper"
)

type Repository struct {
	db *pgxwrapper.Wrapper
}

func New(db *pgxwrapper.Wrapper) *Repository {
	return &Repository{db: db}
}
