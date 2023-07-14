package storage

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqliteStorage(dbPath string) (*Sqlite, error) {
	const op = "storage.sqlite.NewSqliteStorage"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	create table if not exists link(
	    id integer primary key,
	    alias text not null unique,
	    link text not null);
	create index if not exists idx_alias on link(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Sqlite{db: db}, nil
}
