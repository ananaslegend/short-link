package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r PostgresRepository) InsertLink(ctx context.Context, link, alias string) error {
	stmt, err := r.db.Prepare(`
	INSERT INTO link (alias, link)
	VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, alias, link); err != nil {
		if pgErr, ok := lo.ErrorsAs[*pq.Error](err); ok && pgErr.Code == "23505" {
			return ErrAliasAlreadyExists
		}

		return err
	}

	return nil
}
