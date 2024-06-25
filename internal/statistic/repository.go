package statistic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) InsertRows(ctx context.Context, rows Rows) error {
	const op = "internal.statistic.repository.InsertRows"
	stmt, err := r.db.Prepare(
		`insert into statistic (redirect_time_stamp, link, redirect)
				values ($1,$2,$3)`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var errs = make([]error, 0)
	for dimension, metric := range rows {
		if _, err = stmt.ExecContext(ctx, dimension.Timestamp, dimension.Link, metric.Redirect); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", op, err)) // TODO
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
