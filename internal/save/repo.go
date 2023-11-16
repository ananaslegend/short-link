package save

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) InsertLink(ctx context.Context, link, alias string) error {
	const op = "storage.sql.InsertLink"

	stmt, err := r.db.Prepare(`
	insert into link (alias, link)
	values (?,?)
`)
	if err != nil {
		return fmt.Errorf("%r: %w", op, err)
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, alias, link); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return ErrAliasAlreadyExists
		}

		return fmt.Errorf("%r: %w", op, err)
	}

	return nil
}
