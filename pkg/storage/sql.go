package storage

import (
	"database/sql"
	"fmt"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
)

func Close(db *sql.DB, log *slog.Logger) {
	err := db.Close()
	if err != nil {
		log.Error("cant close database", logs.Err(err))
		return
	}
	log.Debug("database closed")
}

func Prepare(db *sql.DB) error {
	const op = "storage.sql.Prepare"

	stmt, err := db.Prepare(`
	create table if not exists link(
	    id integer primary key,
	    alias text not null unique,
	    link text not null);

	create index if not exists idx_alias on link(alias);
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
