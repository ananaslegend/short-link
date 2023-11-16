package redirect

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) SelectLink(ctx context.Context, alias string) (string, error) {
	const op = "storage.sql.SelectLink"

	stmt, err := r.db.Prepare(`
	select link 
	from link
	where alias == ?
`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var link string
	if err = stmt.QueryRowContext(ctx, alias).Scan(&link); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
