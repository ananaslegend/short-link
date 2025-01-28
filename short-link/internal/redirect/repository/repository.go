package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) SelectLink(ctx context.Context, alias string) (string, error) {
	stmt, err := r.db.Prepare(`
	select link 
	from link
	where alias = $1
`)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	var link string
	if err = stmt.QueryRowContext(ctx, alias).Scan(&link); err != nil {
		return "", err
	}

	return link, nil
}
