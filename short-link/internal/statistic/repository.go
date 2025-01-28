package statistic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r PostgresRepository) InsertRows(ctx context.Context, rows Rows) error {
	const op = "internal.statistic.repository.InsertRows"
	stmt, err := r.db.Prepare(
		`insert into statistic (redirect_time_stamp, link, redirect)
				values ($1,$2,$3)`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	var errs = make([]error, 0)
	for dimension, metric := range rows {
		if _, err = stmt.ExecContext(ctx, dimension.Timestamp, dimension.Link, metric.Redirect); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", op, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

type NativeClickHouse struct {
	conn driver.Conn
}

func NewNativeClickHouseRepository(conn driver.Conn) NativeClickHouse {
	return NativeClickHouse{conn: conn}
}

func (r NativeClickHouse) InsertRows(ctx context.Context, rows Rows) error {
	var errs = make([]error, 0)
	for dimension, metric := range rows {
		err := r.conn.AsyncInsert(
			ctx,
			`insert into statistic (redirect_time_stamp, link, redirect)values (?,?,?)`,
			false,
			dimension.Timestamp, dimension.Link, metric.Redirect,
		)

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
