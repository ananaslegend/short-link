package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) InsertLink(ctx context.Context, link, alias string) error {
	const op = "storage.sql.InsertLink"

	stmt, err := r.db.Prepare(`
	insert into link (alias, link)
	values (?,?)
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, alias, link); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return ErrAliasAlreadyExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
