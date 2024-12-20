package repository

import (
	"context"
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLite(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r SQLiteRepository) InsertLink(ctx context.Context, link, alias string) error {
	stmt, err := r.db.Prepare(`
	insert into link (alias, link)
	values (?,?)
`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, alias, link); err != nil {
		if sqliteErr, ok := lo.ErrorsAs[*sqlite3.Error](err); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return ErrAliasAlreadyExists
		}

		return err
	}

	return nil
}
